package gorule

import (
	"github.com/curltech/go-colla-biz/ruleengine"
	"github.com/curltech/go-colla-biz/ruleengine/entity"
	"github.com/curltech/go-colla-biz/ruleengine/service"
	"github.com/curltech/go-colla-core/content"
	entity2 "github.com/curltech/go-colla-core/entity"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

func Fire(packageName string, version string, facts map[string]interface{}) error {
	//1.装载事实
	dataCtx := ast.NewDataContext()
	for k, v := range facts {
		err := dataCtx.Add(k, v)
		if err != nil {
			return err
		}
	}
	knowledgeBase, err := load(packageName, version)
	if err != nil {
		return err
	}
	err = fire(dataCtx, knowledgeBase)
	if err != nil {
		return err
	}

	return nil
}

func getKey(packageName string, version string) string {
	key := "RuleEngine:" + ruleengine.ExecuteType_GoRule + ":" + packageName + ":" + version

	return key
}

func getCacheRule(packageName string, version string) (*ast.KnowledgeBase, bool) {
	v, ok := ruleengine.MemCache.Get(getKey(packageName, version))
	if ok {
		return v.(*ast.KnowledgeBase), true
	}

	return nil, false
}

func setCacheRule(packageName string, version string, rule *ast.KnowledgeBase) {
	ruleengine.MemCache.SetDefault(getKey(packageName, version), rule)
}

/**
	rule SpeedUp "When testcar is speeding up we keep increase the speed." salience 10  {
    when
        TestCar.SpeedUp == true && TestCar.Speed < TestCar.MaxSpeed
    then
        TestCar.Speed = TestCar.Speed + TestCar.SpeedIncrement;
        DistanceRecord.TotalDistance = DistanceRecord.TotalDistance + TestCar.Speed;
}
*/
func loadResources(packageName string, version string) []pkg.Resource {
	//2.创建知识库
	svc := service.GetRuleDefinitionService()
	ruleDefinitions := make([]*entity.RuleDefinition, 0)
	ruleDefinition := entity.RuleDefinition{}
	ruleDefinition.ExecuteType = ruleengine.ExecuteType_GoRule
	ruleDefinition.Status = entity2.EntityStatus_Effective
	ruleDefinition.PackageName = packageName
	ruleDefinition.Version = version
	svc.Find(&ruleDefinitions, &ruleDefinition, "", "")
	var err error
	resources := make([]pkg.Resource, 0)
	for _, v := range ruleDefinitions {
		raw := v.RawContent
		var drl []byte
		if raw == "" {
			drl, err = content.FileContent.Read(v.ContentId)
			if err != nil {
				continue
			}
		} else {
			drl = []byte(raw)
		}
		resource := pkg.NewBytesResource(drl)
		resources = append(resources, resource)
	}

	return resources
}

func load(packageName string, version string) (*ast.KnowledgeBase, error) {
	//1.判断是否已经存在
	knowledgeBase, ok := getCacheRule(packageName, version)
	if ok {
		return knowledgeBase, nil
	}
	resources := loadResources(packageName, version)
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	//3.加载规则，从字节流，文件，字符串，数据库等来源
	//	drls := `
	//	rule CheckValues "Check the default values" salience 10 {
	//    when
	//        MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
	//    then
	//        MF.WhatToSay = MF.GetWhatToSay("Hello Grule");
	//        Retract("CheckValues");
	//}
	//`
	err := ruleBuilder.BuildRuleFromResources(packageName, version, resources)
	if err != nil {
		return nil, err
	}
	knowledgeBase = knowledgeLibrary.NewKnowledgeBaseInstance(packageName, version)
	setCacheRule(packageName, version, knowledgeBase)

	return knowledgeBase, nil
}

func fire(dataCtx ast.IDataContext, knowledgeBase *ast.KnowledgeBase) error {
	//3.事实与知识库结合执行
	engine := engine.NewGruleEngine()
	err := engine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		return err
	}

	return nil
}
