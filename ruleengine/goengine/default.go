package goengine

import (
	"errors"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"github.com/curltech/go-colla-biz/ruleengine"
	"github.com/curltech/go-colla-biz/ruleengine/entity"
	"github.com/curltech/go-colla-biz/ruleengine/service"
	"github.com/curltech/go-colla-core/content"
	entity2 "github.com/curltech/go-colla-core/entity"
)

func Fire(packageName string, version string, facts map[string]interface{}) error {
	//1.装载事实
	/**
	不要注入除函数和结构体指针以外的其他类型(如变量)
	*/
	dataCtx := context.NewDataContext()
	for k, v := range facts {
		dataCtx.Add(k, v)
	}
	ruleBuilder, err := load(packageName, version, dataCtx)
	if err != nil {
		return err
	}
	err = fire(ruleBuilder)
	if err != nil {
		return err
	}

	return nil
}

func getKey(packageName string, version string) string {
	key := "RuleEngine:" + ruleengine.ExecuteType_GoEngine + ":" + packageName + ":" + version

	return key
}

func getCacheRule(packageName string, version string) (string, bool) {
	v, ok := ruleengine.MemCache.Get(getKey(packageName, version))
	if ok {
		return v.(string), true
	}

	return "", false
}

func setCacheRule(packageName string, version string, rule string) {
	ruleengine.MemCache.SetDefault(getKey(packageName, version), rule)
}

func load(packageName string, version string, dataCtx *context.DataContext) (*builder.RuleBuilder, error) {
	rule, ok := getCacheRule(packageName, version)
	if !ok {
		//2.创建知识库
		svc := service.GetRuleDefinitionService()
		ruleDefinitions := make([]*entity.RuleDefinition, 0)
		ruleDefinition := entity.RuleDefinition{}
		ruleDefinition.ExecuteType = ruleengine.ExecuteType_GoEngine
		ruleDefinition.Status = entity2.EntityStatus_Effective
		ruleDefinition.PackageName = packageName
		ruleDefinition.Version = version
		svc.Find(&ruleDefinitions, &ruleDefinition, "", 0, 0, "")
		if len(ruleDefinitions) > 0 {
			for _, ruleDefinition := range ruleDefinitions {
				var drl = ruleDefinition.RawContent
				if drl == "" {
					v, err := content.FileContent.Read(ruleDefinition.ContentId)
					if err != nil {
						continue
					} else {
						drl = string(v)
					}
				}
				if drl != "" {
					rule = rule + "\n\n" + drl
				}
			}
		}
	}

	//读取规则，多个规则可以回车符拼在一起
	//初始化规则引擎
	if rule != "" {
		ruleBuilder := builder.NewRuleBuilder(dataCtx)
		err := ruleBuilder.BuildRuleFromString(rule)
		if err == nil {
			setCacheRule(packageName, version, rule)

			return ruleBuilder, nil
		} else {
			return nil, err
		}
	}
	return nil, errors.New("NoRule")
}

func fire(ruleBuilder *builder.RuleBuilder) error {
	//3.事实与知识库结合执行
	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	err := eng.Execute(ruleBuilder, true)
	if err != nil {
		return err
	}

	return nil
}
