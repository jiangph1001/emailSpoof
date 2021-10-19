package main

import (
	"bufio"
	"email/smtp"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

var config Config
var logger *log.Logger

func getResponse(reader *bufio.Reader) string {
	response := ""
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		response += line
		if line[3] == ' ' {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return response
}
func getSingleResponse(reader *bufio.Reader) string {
	res, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(res))
	time.Sleep(300 * time.Millisecond)
	return string(res)
}

func GetData(mail Mail) string {
	from := mail.Data.From
	to := mail.Data.To
	subject := mail.Data.Subject
	xffStr := string(0xff)
	data := fmt.Sprintf("From%s:%s\nFrom:%s\nTo:%s\nSubject:%s\r\n.\r\n", xffStr, mail.Header.MailFrom, from, to, subject)
	return data
}
func GetCommand(mail Mail) []string {
	var commandList []string
	ehloCommand := fmt.Sprintf("EHLO %s\r\n", mail.Header.Ehlo)
	mailFromCommand := fmt.Sprintf("MAIL FROM:<%s>\r\n", mail.Header.MailFrom)
	rcptToCommand := fmt.Sprintf("RCPT TO:<%s>\r\n", mail.Header.RcptTo)
	commandList = append(commandList, ehloCommand)
	//if mail.Login == 1 {
	//	commandList = append(commandList, "AUTH LOGIN")
	//	commandList = append(commandList, base64.StdEncoding.EncodeToString([]byte(mail.Auth.User)))
	//	commandList = append(commandList, base64.StdEncoding.EncodeToString([]byte(mail.Auth.Password)))
	//}
	commandList = append(commandList, mailFromCommand)
	commandList = append(commandList, rcptToCommand)
	commandList = append(commandList, "DATA\r\n")
	commandList = append(commandList, GetData(mail))
	commandList = append(commandList, "QUIT\r\n")
	return commandList
}
func getResCode(response string) int {
	resCode, _ := strconv.ParseInt(response[:3], 0, 32)
	return int(resCode)
}
func SendCommand(conn net.Conn, command string) string {
	reader := bufio.NewReader(conn)
	writeLog(logger, "> "+command, config.SendLogScreen)
	fmt.Fprintf(conn, command)
	//response := getSingleResponse(reader)
	response := getResponse(reader)
	writeLog(logger, "< "+response, config.ReceiveLogScreen)
	return response
}
func sendEmailMine(mail Mail) {
	commandList := GetCommand(mail)
	address := fmt.Sprintf("%s:%d", mail.Host, mail.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
		logger.Fatal("connect", address, "fail", err)
	}
	defer conn.Close()
	logger.Printf("connect %s success", address)
	reader := bufio.NewReader(conn)
	getResponse(reader)
	for _, command := range commandList {
		SendCommand(conn, command)
	}
}

func sendEmailWithAuth(mail Mail) error {
	auth := smtp.PlainAuth("", mail.Auth.User, mail.Auth.Password, mail.Host)
	msg := []byte(GetData(mail))
	to := []string{mail.Header.RcptTo}
	address := fmt.Sprintf("%s:%d", mail.Host, mail.Port)
	err := smtp.SendMail(address, auth, mail.Header.MailFrom, to, msg)
	return err
}

func main() {
	logger = getLogger()
	config = readConfigJsonFile()
	mailInfo := readEmailJsonFile()
	for cnt, mail := range mailInfo.Mails {
		if mail.Send == 0 {
			continue
		}
		log.Print("send email ", cnt)
		if mail.Login == 1 {
			err := sendEmailWithAuth(mail)
			if err != nil {
				log.Print("err:", err)
			}
		} else {
			sendEmailMine(mail)
		}
	}
}
