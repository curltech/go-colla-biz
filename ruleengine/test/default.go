package test

import (
	"fmt"
	"github.com/curltech/go-colla-biz/ruleengine"
	"github.com/curltech/go-colla-biz/ruleengine/entity"
	"github.com/curltech/go-colla-biz/ruleengine/goengine"
	"github.com/curltech/go-colla-biz/ruleengine/service"
	entity2 "github.com/curltech/go-colla-core/entity"
	"github.com/kataras/golog"
)

const (
	rule = `
rule "测试" "测试描述"  salience 0 
begin
		// 重命名函数 测试
		Sout("XXXXXXXXXX")
		// 普通函数 测试
		Hello()
		//结构提方法 测试
		User.Say()
		// if
		if 7 == User.GetNum(7){
			//自定义变量 和 加法 测试
			variable = "hello" + " world"
			// 加法 与 内建函数 测试
			User.Name = "hhh" + strconv.FormatBool(true)
			//结构体属性、方法调用 和 除法 测试
			User.Age = User.GetNum(89767999999) / 10000000
			//布尔值设置 测试
			User.Male = false
			//规则内自定义变量调用 测试
			User.Print(variable)
			//float测试	也支持科学计数法		
			f = 9.56			
			PrintReal(f)
			//嵌套if-else测试
			if false	{
				Sout("嵌套if测试")
			}else{
				Sout("嵌套else测试")
			}
		}else{ //else
			//字符串设置 测试
			User.Name = "yyyy"
		}
end`
)

type User struct {
	Name string
	Age  int
	Male bool
}

func (u *User) GetNum(i int64) int64 {
	return i
}

func (u *User) Print(s string) {
	fmt.Println(s)
}

func (u *User) Say() {
	fmt.Println("hello world")
}

func exe(packageName string, version string, facts map[string]interface{}) {
	err := goengine.Fire(packageName, version, facts)
	if err != nil {
		golog.Errorf("execute rule error: %v", err)
	}
}

func Hello() {
	fmt.Println("hello")
}

func PrintReal(real float64) {
	fmt.Println(real)
}

func Test() {
	user := &User{
		Name: "Calo",
		Age:  0,
		Male: true,
	}
	/**
	不要注入除函数和结构体指针以外的其他类型(如变量)
	*/
	facts := make(map[string]interface{}, 0)
	//注入结构体指针
	facts["User"] = user
	//重命名函数,并注入
	facts["Sout"] = fmt.Println
	//直接注入函数
	facts["Hello"] = Hello
	facts["PrintReal"] = PrintReal
	exe("test", "1.0.0", facts)

	golog.Infof("user.Age=%d,Name=%s,Male=%t", user.Age, user.Name, user.Male)
}

func Prepare() {
	svc := service.GetRuleDefinitionService()
	ruleDefinition := entity.RuleDefinition{}
	ruleDefinition.ExecuteType = ruleengine.ExecuteType_GoEngine
	ruleDefinition.Status = entity2.EntityStatus_Effective
	ruleDefinition.PackageName = "test"
	ruleDefinition.Version = "1.0.0"
	ruleDefinition.RawContent = rule
	svc.Insert(&ruleDefinition)
}
