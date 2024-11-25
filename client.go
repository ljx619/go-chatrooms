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
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
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
	select {}
}
