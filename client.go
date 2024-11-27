package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//连接server
	conn, err := net.Dial("tcp", serverIp+":"+strconv.Itoa(client.ServerPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn

	//返回对象
	return client
}

// SelectUsers 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write err:", err)
		return
	}
}

// PrivateChat 私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>请输入聊天对象[用户名],exit退出:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>请输入消息内容,exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("client.conn.Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>请输入消息内容,exit退出:")
			fmt.Scanln(&chatMsg)

		}

		//remoteName = ""
		client.SelectUsers()
		fmt.Println(">>>请输入聊天对象[用户名],exit退出:")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	var chatMsg string

	//提醒用户输入内容
	fmt.Println(">>>请输入聊天内容,exit退出:")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		//发给服务器
		//消息不为空咋发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("client.conn.Write err:", err)
			}
		}

		//循环处理
		chatMsg = ""
		fmt.Println(">>>请输入聊天内容,exit退出:")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>请输入合法范围内的数字<<<")
		return false
	}
}

func (client *Client) run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			//公聊模式
			client.PublicChat()
			break
		case 2:
			//私聊模式
			client.PrivateChat()
			break
		case 3:
			//更新用户名
			client.UpdateName()
			break
		}
	}
}

// DealResponse 处理server回传过来的消息 这里显示到标准输入输出
func (client *Client) DealResponse() {
	//一旦client.conn有数据就直接拷贝到 stdout 输出中 永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

var serverIp string
var serverPort int

//./client -ip 127.0.0.1 -port 8080

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址,默认是127.0.0.1")
	flag.IntVar(&serverPort, "port", 8080, "设置服务器端口,默认是8080")
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>链接服务器失败...")
		return
	}

	//单独开启一个goroutine处理server回传的消息
	go client.DealResponse()

	fmt.Println(">>>链接服务器成功...")

	//启动客户端的业务
	client.run()
}
