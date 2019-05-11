package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var email = ""         //TomatoPay邮箱
var password = ""      //TomatoPay密码
var threshold = 100.00 //提现阈值，单位：元

func main() {
	const urlHome = "https://b.fanqieui.com/"

	client := &http.Client{}

	//登录
	requestBody := url.Values{}
	requestBody.Set("email", email)
	requestBody.Set("password", password)
	request, err := http.NewRequest("POST", urlHome+"_login.php", strings.NewReader(requestBody.Encode()))
	panicError(err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	panicError(err)
	cookies := response.Cookies()
	body, err := ioutil.ReadAll(response.Body)
	panicError(err)
	resultObject := struct {
		Code    string `json:"code"`
		Message string `json:"msg"`
	}{}
	panicError(json.Unmarshal(body, &resultObject))
	if resultObject.Code != "1" {
		panicLog(resultObject.Message)
	}
	panicError(response.Body.Close())

	//准备提现数据
	request, err = http.NewRequest("GET", urlHome+"dashboard/withdrawal.php", nil)
	panicError(err)
	for _, cookieIndex := range cookies {
		request.AddCookie(cookieIndex)
	}
	response, err = client.Do(request)
	panicError(err)
	body, err = ioutil.ReadAll(response.Body)
	panicError(err)
	head := `<input hidden type="text" id="token" name="token" value="`
	token := regexFind(body, head)

	head = `<input class="form-control" type="text" id="cny" name="cny" placeholder="您可以提现¥ `
	balance := regexFind(body, head)
	balanceFloat, err := strconv.ParseFloat(balance, 64)
	panicError(err)
	if balanceFloat < threshold {
		panicLog("余额为：" + balance + "，未达到阈值")
	}
	for _, ck := range response.Cookies() {
		cookies = append(cookies, ck)
	}
	panicError(response.Body.Close())

	//申请提现
	requestBody = url.Values{}
	requestBody.Set("token", token)
	requestBody.Set("cny", balance)
	request, err = http.NewRequest("POST", urlHome+"dashboard/_withdrawal.php", strings.NewReader(requestBody.Encode()))
	panicError(err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Referer", urlHome+"dashboard/withdrawal.php")
	for _, cookieIndex := range cookies {
		request.AddCookie(cookieIndex)
	}
	response, err = client.Do(request)
	panicError(err)
	body, err = ioutil.ReadAll(response.Body)
	panicError(err)
	panicError(json.Unmarshal(body, &resultObject))
	if resultObject.Code != "1" {
		panicLog(resultObject.Message)
	}

	println("已申请提现" + balance + "人民币")
}

func regexFind(context []byte, head string) string {
	const tail = `">`
	compiled, err := regexp.Compile(strings.ReplaceAll(head, ` `, `\s`) + `.*?` + tail)
	panicError(err)
	match := string(compiled.Find(context))
	return match[len(head) : len(match)-len(tail)]
}

var allLog = time.Now().String() + ":"

func panicError(err error) {
	if err != nil {
		panicLog(err.Error())
	}
}

func panicLog(log string) {
	allLog += "\n" + log + "\n"
	logFile, _ := os.OpenFile("AutoTomatoPay.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.ModeAppend)
	_, _ = logFile.WriteString(allLog)
	panic(allLog)
}
