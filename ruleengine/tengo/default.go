package tengo

import (
	"errors"
	"github.com/curltech/go-colla-biz/ruleengine"
	"github.com/curltech/go-colla-biz/ruleengine/entity"
	"github.com/curltech/go-colla-biz/ruleengine/service"
	"github.com/curltech/go-colla-core/content"
	entity2 "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/logger"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

func Fire(packageName string, version string, facts map[string]interface{}) error {
	//1.装载事实
	compileds, err := load(packageName, version)
	if err != nil {
		return err
	}
	err = fire(compileds, facts)
	if err != nil {
		return err
	}

	return nil
}

func getKey(packageName string, version string) string {
	key := "RuleEngine:" + ruleengine.ExecuteType_TenGo + ":" + packageName + ":" + version

	return key
}

func getCacheRule(packageName string, version string) ([]*tengo.Compiled, bool) {
	v, ok := ruleengine.MemCache.Get(getKey(packageName, version))
	if ok {
		return v.([]*tengo.Compiled), true
	}

	return nil, false
}

func setCacheRule(packageName string, version string, rule []*tengo.Compiled) {
	ruleengine.MemCache.SetDefault(getKey(packageName, version), rule)
}

func loadResources(packageName string, version string) []*tengo.Script {
	//2.创建知识库
	svc := service.GetRuleDefinitionService()
	ruleDefinitions := make([]*entity.RuleDefinition, 0)
	ruleDefinition := entity.RuleDefinition{}
	ruleDefinition.ExecuteType = ruleengine.ExecuteType_TenGo
	ruleDefinition.Status = entity2.EntityStatus_Effective
	ruleDefinition.PackageName = packageName
	ruleDefinition.Version = version
	svc.Find(&ruleDefinitions, &ruleDefinition, "", 0, 0, "")
	var err error
	resources := make([]*tengo.Script, 0)
	for _, v := range ruleDefinitions {
		raw := v.RawContent
		var src []byte
		if raw == "" {
			src, err = content.FileContent.Read(v.ContentId)
			if err != nil {
				continue
			}
		} else {
			src = []byte(raw)
		}
		script := tengo.NewScript(src)
		resources = append(resources, script)
	}

	return resources
}

func load(packageName string, version string) ([]*tengo.Compiled, error) {
	//1.判断是否已经存在
	compileds, ok := getCacheRule(packageName, version)
	if ok {
		return compileds, nil
	}
	scripts := loadResources(packageName, version)
	if scripts != nil && len(scripts) > 0 {
		compileds = make([]*tengo.Compiled, 0)
		for _, script := range scripts {
			compiled, err := script.Compile()
			if err != nil {
				logger.Errorf(err.Error())
				continue
			}
			script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
			//mods := tengo.NewModuleMap()
			//mods.AddSourceModule("double", []byte(`export func(x) { return x * 2 }`))
			//script.SetImports(mods)
			compileds = append(compileds, compiled)
		}
		setCacheRule(packageName, version, compileds)
	} else {
		return nil, errors.New("NoScript")
	}

	return compileds, nil
}

func fire(compileds []*tengo.Compiled, facts map[string]interface{}) error {
	//3.事实与知识库结合执行
	var err error
	if compileds != nil && len(compileds) > 0 {
		for _, compiled := range compileds {
			for k, v := range facts {
				_ = compiled.Set(k, v)
			}
			er := compiled.Run()
			if er != nil {
				err = er
				logger.Errorf(err.Error())
				continue
			}
			//compiled.GetAll()
		}
	}

	return err
}
