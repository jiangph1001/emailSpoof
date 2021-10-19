package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type MailList struct {
	Mails []Mail `json:"mails"`
}
type Mail struct {
	Send   int            `json:"send"`
	Login  int            `json:"login"`
	Host   string         `json:"host"`
	Port   int            `json:"port"`
	Header HeaderStruct   `json:"header"`
	Data   DataStruct     `json:"data"`
	Auth   Authentication `json:"auth"`
}
type HeaderStruct struct {
	Ehlo     string `json:"ehlo"`
	MailFrom string `json:"mail from"`
	RcptTo   string `json:"rcpt to"`
}
type DataStruct struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
}
type Config struct {
	EmailConfig      string `json:"email_config"`
	SendLogScreen    int    `json:"send_log_screen"`
	ReceiveLogScreen int    `json:"receive_log_screen"`
}
type Authentication struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func readEmailJsonFile() MailList {
	jsonContent, err := ioutil.ReadFile("email.json")
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
	}
	mails := MailList{}
	err = json.Unmarshal(jsonContent, &mails)
	if err != nil {
		fmt.Println("解析数据失败", err)
	}
	return mails
}
func readConfigJsonFile() Config {
	jsonContent, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
	}
	config := Config{}
	err = json.Unmarshal(jsonContent, &config)
	if err != nil {
		fmt.Println("解析数据失败", err)
	}
	return config
}
