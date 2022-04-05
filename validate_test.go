package easyValidator

import (
	"fmt"
	"net/http"
	"testing"
)

type Home struct {
	Addr string `form:"addr" check:"max=3"`
	Xxx int `form:"xxx"`
}

type Student struct {
	Home
	Age byte `form:"age" check:"gt=18"`
	Height uint8 `form:"height" check:"select=100 200,gt=50"`
	Name string `form:"name" check:"max=4"`
}



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
	validator := NewStructValidator()
	err := validator.ValidateStruct(a)
	if err != nil {
		t.Log(err.Error())
	}
}

func TestHttpReqValidator(t *testing.T) {
	stu:= Student{}
	req,_ := http.NewRequest("GET", "http://my.url.com/root?age=19&name=jack&addr=123", nil)
	validator := NewHttpRequestValidator()
	if err := validator.BindHttpRequest(req, &stu); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n",stu)
}