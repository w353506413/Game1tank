package main

import (
	"bufio"
	"net"
)

func main() {
	listenner, _ := net.Listen("tcp", ":1024")
	println("启动服务器....")
	for {
		conn, _ := listenner.Accept() // 持续监听客户端连接
		go ClientLogic(conn)
	}
}
func ClientLogic(conn net.Conn) {
	//从客户端接受数据
	s, _ := bufio.NewReader(conn).ReadString('\n')
	println("有客户端发来的消息:", s)

	//发送消息给客户端
	conn.Write([]byte("东东你好\n"))

	//关闭连接
	conn.Close()

}
