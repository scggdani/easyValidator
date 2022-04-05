## 参数校验器--easyValidator
> 轻量级的`go-playground/validator` 

### 1. 功能与特性
- 能够对结构体字段进行校验
- 能够对http请求的参数进行校验，并将参数解析到对应的结构体
- 支持结构体嵌套

### 2. 用法
- 标签类型

标签`form` 为参数名标签，会根据此标签的值从http请求中解析出相应字段。
标签`check` 为校验标签，会根据改标签里的属性对字段进行校验。


相关属性：  
```go
"it" // 小于, 针对数字类型
"gt" // 大于, 针对数字类型
"ite" // 小于等于, 针对数字类型
"gte" // 大于等于, 针对数字类型
"max" // 最大长度, 针对字符串、切片
"min" // 最小长度, 针对字符串、切片
"len" // 指定长度, 针对字符串、切片 
"select" // 应该为其中之一, 针对字符串和数值
```  
- 结构体定义
> - 校验结构体时字段大小写不限
> - 若需要从http请求中解析参数，请将传入结构体中待解析的字段首字母大写
> - 结构体嵌套时不得出现空指针，否则会报错
```go
type Home struct {
	Addr string `form:"addr" check:"max=3"` // max=3表示限定该字段长度最大为3
	Xxx int `form:"xxx"`
}

type Student struct {
	Home
	Age byte `form:"age" check:"gt=18"` //gt=18表示该字段值需要大于18
	Height uint8 `form:"height" check:"select=100 200,gt=50"` // select=100 200表示该字段值为100,200中的一个
	Name string `form:"name" check:"max=4"`
}
```

- 校验结构体
```go
func TestStructValidator(t *testing.T) {
	a := Student{
		Home: Home{
			Addr: "1234",
			Xxx:  0,
		},
		Age:    22,
		Height: 100,
		Name:   "123",
	}
	validator := NewStructValidator() // 创建结构体校验器
	err := validator.ValidateStruct(a) // 调用ValidateStruct方法进行校验
	if err != nil {
		t.Log(err.Error())
	}
}
```
结果：`字段[addr]错误, 值为[1234], 校验规则为[最大长度为3]` 

- http请求解析
```go
func TestHttpReqValidator(t *testing.T) {
    stu:= Student{} // 进行http解析时，嵌套结构体中不允许出现空指针
    req,_ := http.NewRequest("GET", "http://my.url.com/root?age=19&name=jack&addr=123", nil)
    validator := NewHttpRequestValidator() // 创建http参数校验解析器
    if err := validator.BindHttpRequest(req, &stu); err != nil { // err == nil时将参数成功的解析到stu中
        fmt.Println(err)
    }
    fmt.Printf("%+v\n",stu)
}
```
结果：`{Home:{Addr:123 Xxx:0} Age:19 Height:0 Name:jack}`  

### 下一步计划：
- 支持从上下文Context中解析参数并校验