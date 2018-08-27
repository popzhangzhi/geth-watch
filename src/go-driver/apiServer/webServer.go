package apiServer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WebServerBase() {
	fmt.Println("this is a WebServerBase")
	//注册路由
	http.HandleFunc("/login", loginTask)
	//开启监听服务
	err := http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		fmt.Println("ListenAndServe err:", err.Error())
	}
}

func loginTask(w http.ResponseWriter, req *http.Request) {

	fmt.Println("loginTask is running")

	//time.Sleep(time.Second*2)
	//req.URL.Query()
	req.ParseForm()
	userName_arr, userNameErr := req.Form["userName"]
	pwd_arr, pwdErr := req.Form["pwd"]
	if !(userNameErr && pwdErr) {
		fmt.Println("输入不合法")
		return
	}
	userName := userName_arr[0]
	pwd := pwd_arr[0]
	result := NewJsonResponse()

	fmt.Println(userName, pwd)

	if userName == "1" && pwd == "2" {
		result.Code = 0
		data := rel{}
		data.Data = []dataOutput{dataOutput{"a", "b", 12}, dataOutput{"1", "2", 3}}
		data.Info = "string"

		result.Data = data
		result.Msg = "登录成功"

	} else {
		result.Code = 1
		//result.Data=nil
		result.Msg = "账号密码不正确"
	}
	fmt.Println(result)
	bytes, _ := json.Marshal(result)
	fmt.Fprint(w, string(bytes))

}

type JsonResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewJsonResponse() *JsonResponse {
	return &JsonResponse{}
}

type rel struct {
	Data []dataOutput `json:"data"`
	Info string       `json:"info"`
}

type dataOutput struct {
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
	Age      int8   `json:"age"`
}
