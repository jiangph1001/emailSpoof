package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func getLogger() *log.Logger {
	currentTime := time.Now()
	t1 := currentTime.Year()   //年
	t2 := currentTime.Month()  //月
	t3 := currentTime.Day()    //日
	t4 := currentTime.Hour()   //小时
	t5 := currentTime.Minute() //分钟
	logName := fmt.Sprintf("log/%d-%d-%d-%d-%d.log", t1, t2, t3, t4, t5)
	file, err := os.OpenFile("log/smtp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("fail to create %s file:%s\n", logName, err)
	}
	logger := log.New(file, "", log.LstdFlags)
	return logger
}

func writeLog(logger *log.Logger, record string, flag int) {
	logger.Print(record)
	if flag == 1 {
		log.Print(record)
	}
}
