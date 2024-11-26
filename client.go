package main

import (
	"flag"
	"fmt"
	"net"
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

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.私聊模式")
	fmt.Println("2.公聊模式")
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
			break
		case 2:
			//私聊模式
			break
		case 3:
			//更新用户名
			break
		}
	}
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

	fmt.Println(">>>链接服务器成功...")

	//启动客户端的业务
	client.run()
}
