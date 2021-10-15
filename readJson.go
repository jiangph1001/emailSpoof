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
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Header HeaderStruct `json:"header"`
	Data DataStruct `json:"data"`
}

type HeaderStruct struct {
	Ehlo   string `json:"ehlo"`
	MailFrom string `json:"mail from"`
	RcptTo   string `json:"rcpt to"`
}
type DataStruct struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
}
func readJsonFile() MailList{
	jsonContent, err := ioutil.ReadFile("config.json")
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

