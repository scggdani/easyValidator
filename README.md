## 参数校验器--easyValidator
> 一款轻量级的参数校验与解析器，与`go-playground/validator`相比，信息提示更人性化、可自由集成到各种项目中 

### 1. 功能与特性
- 能够对结构体字段进行校验
- 能够对http请求的参数进行校验，并将参数解析到对应的结构体
- 能将Context上下文中保存的值解析到结构体并进行校验  
- 支持结构体嵌套

### 2. 用法
- 标签类型

标签`form` 为参数名标签，会根据此标签的值从http请求中解析出相应字段。

标签`check` 为校验标签，会根据改标签里的属性对字段进行校验。

标签`default` 为默认值标签，当找不到对应字段值的时候会将改标签定义的值赋给对应字段。


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
	Xxx int `form:"xxx" default:"666"` // default=666表示当解析不出该字段值时使用666作为默认值
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

```go
func TestStructValidator_Select(t *testing.T) {
	a := Student{
		Home: Home{
			Addr: "123",
			Xxx:  0,
		},
		Age:    22,
		Height: 150,
		Name:   "123",
	}
	validator := NewStructValidator()
	err := validator.ValidateStruct(a)
	if err != nil {
		t.Log(err.Error())
	}
}
```
结果：`字段[Height]错误, 值为[150], 校验规则为[值为下列之一:{100,200};大于50]` 

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
结果：`{Home:{Addr:123 Xxx:666} Age:19 Height:0 Name:jack}`  

- Context上下文解析
```go
func TestContextValidator_TypeError(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "addr", "123")
	ctx = context.WithValue(ctx, "age", int(22)) // 当context中指定键的值的类型与对应结构体字段类型不符时会报错
	ctx = context.WithValue(ctx, "height", uint(100))
	a := Student{}
	validator := NewContextValidator()
	if err := validator.BindContext(ctx, &a); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v\n", a)
}
```
由于`context` 中`age` 为`int` 类型，而所传结构体中`age` 字段为`byte` 类型，故会报错，结果为：
```go
context解析出现类型错误,结构体中类型为[uint8],context中对应值的类型为[int]
{Home:{Addr:123 Xxx:666} Age:0 Height:0 Name:}
```

正常情况下：
```go
func TestContextValidator(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "addr", "123")
	ctx = context.WithValue(ctx, "height", uint(100))
	a := Student{}
	validator := NewContextValidator()
	if err := validator.BindContext(ctx, &a); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v\n", a)
}
```
结果为：`{Home:{Addr:123 Xxx:666} Age:22 Height:100 Name:}` 

### 下一步计划：
- 算法优化
- 更多的校验属性支持