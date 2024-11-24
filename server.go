package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// NewServer 创建一个 Server 的对外函数
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// ListenMessage 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		//将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// BroadCast 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + user.Name + "]" + msg
	this.Message <- sendMsg
}

func (this *Server) Handle(conn net.Conn) {
	fmt.Println("成功建立链接")

	user := NewUser(conn, this)

	//广播当前用户上线消息
	user.Online()

	// 监听是否有不活跃的 channel
	isLive := make(chan bool)

	//接收客户端传递过来的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户的消息(去除\n回车符)
			msg := string(buf[:n-1])

			//用户针对msg进行处理
			user.DoMessage(msg)

			//用户的任意消息，都会重置自己的活跃度
			isLive <- true
		}
	}()
	// 当前 handel 阻塞
	for {
		select {
		case <-isLive:
			//不做任何事情 为了重置下面的 计时器
			//这里有一个技巧 case 没用 break 所以执行到这里之后也会执行下面的
		case <-time.After(time.Second * 30):
			//30s后超时
			//将 当前 user 踢出去
			user.SendMsg("已超时,被踢出")

			//销毁用户资源
			close(user.C)
			conn.Close()

			//退出当前的Handler
			return
		}
	}

}

// Start 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", this.Ip+":"+strconv.Itoa(this.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	//启动监听 Message的 goroutine
	go this.ListenMessage()

	for {

		//accept
		accept, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		//do handler
		go this.Handle(accept)

	}

}
