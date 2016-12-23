package src

import (
	"net/url"
	"bytes"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"sort"
	"crypto/md5"
	"io"
	"net/http"
)

//发送邮件
func HttpSendMail(to string, body string) {
	RequestURI := "http://sendcloud.sohu.com/webapi/mail.send.json"
	//不同于登录SendCloud站点的帐号，您需要登录后台创建发信子帐号，使用子帐号和密码才可以进行邮件的发送。
	PostParams := url.Values{
		"api_user": {"usen"},
		"api_key":  {"key"},
		"from":     {"service@sendcloud.im"},
		"fromname": {"xxxx"},
		"to":       {to},
		"subject":  {"激活邮件"},
		"html":     {body},
	}
	PostBody := bytes.NewBufferString(PostParams.Encode())
	ResponseHandler, err := http.Post(RequestURI, "application/x-www-form-urlencoded", PostBody)
	if err != nil {
		panic(err)
	}
	defer ResponseHandler.Body.Close()
	BodyByte, _ := ioutil.ReadAll(ResponseHandler.Body)
	fmt.Println(string(BodyByte))
}

//发送短信
func HttpSendSMS(phone string, code string) {
	sms_user := "user"
	sms_key := "key"
	sms_url := "http://sendcloud.sohu.com/smsapi/send"
	vars := map[string]string{
		`%Code%`: code,
	}
	jsonVars, _ := json.Marshal(&vars)
	PostParams := url.Values{
		`smsUser`:    {sms_user},
		`templateId`: {`3623`},
		`msgType`:    {`0`},
		`phone`:      {phone},
		`vars`:       {string(jsonVars)},
	}
	paramsKeyS := make([]string,0,len(PostParams))
	for k, _ := range PostParams {
		paramsKeyS = append(paramsKeyS, k)
	}
	sort.Strings(paramsKeyS)
	sb := sms_key + "&"
	for _, v := range paramsKeyS {
		sb += fmt.Sprintf("%s=%s&", v, PostParams.Get(v))
	}
	sb += sms_key
	hashMd5 := md5.New()
	io.WriteString(hashMd5, sb)
	sign := fmt.Sprintf("%x", hashMd5.Sum(nil))
	PostParams.Add("signature",sign)

	PostBody := bytes.NewBufferString(PostParams.Encode())
	ResponseHandler, _ := http.Post(sms_url, "application/x-www-form-urlencoded", PostBody)
	defer ResponseHandler.Body.Close()
	BodyByte, _ := ioutil.ReadAll(ResponseHandler.Body)
	fmt.Println(string(BodyByte))

}