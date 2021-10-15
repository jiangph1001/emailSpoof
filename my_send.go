package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func get_response(reader *bufio.Reader) {
	for {
		res, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(res))
	}
}
func getSingleResponse(reader *bufio.Reader) {
	res, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))
	time.Sleep(time.Duration(1) * time.Second)
}
func SendCommand(conn net.Conn, str string) {
	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, str)
	getSingleResponse(reader)
}
func GetData(mailData DataStruct) string {
	from := mailData.From
	to := mailData.To
	subject := mailData.Subject
	data := fmt.Sprintf("From:%s\nTo:%s\nSubject:%s\r\n.\r\n", from, to, subject)
	return data
}

func GetCommand(mail Mail) []string {
	var commandList []string
	ehloCommand := fmt.Sprintf("EHLO %s\r\n", mail.Header.Ehlo)
	mailFromCommand := fmt.Sprintf("MAIL FROM:<%s>\r\n", mail.Header.MailFrom)
	rcptToCommand := fmt.Sprintf("RCPT TO:<%s>\r\n", mail.Header.RcptTo)

	commandList = append(commandList, ehloCommand)
	commandList = append(commandList, mailFromCommand)
	commandList = append(commandList, rcptToCommand)
	commandList = append(commandList, "DATA\r\n")
	commandList = append(commandList, GetData(mail.Data))
	commandList = append(commandList, "QUIT\r\n")
	return commandList
}
func sendEmail(mail Mail) {
	commandList := GetCommand(mail)
	address := fmt.Sprintf("%s:%d", mail.Host, mail.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for _, v := range commandList {
		fmt.Printf("%s", v)
		SendCommand(conn, v)
	}
}
func main() {
	mailInfo := readJsonFile()
	//mailNum = len(mailInfo.Mails)
	for _, mail := range mailInfo.Mails {
		sendEmail(mail)
	}
}
