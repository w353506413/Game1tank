package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

func handleConnection(conn net.Conn, talkChan map[int]chan string) {
	//fmt.Printf("%p\n", talkChan)  //用以检查是否是传过来的指针

	/*
		定义当前用户的uid
	*/
	var curUid int

	var err error

	/*
		定义关闭通道
	*/
	var closed = make(chan bool)

	defer func() {
		fmt.Println("defer do : conn closed")
		conn.Close()
		fmt.Printf("delete userid [%v] from talkChan", curUid)
		delete(talkChan, curUid)
	}()

	/**
	 * 提示用户设置自己的uid， 如果没设置，则不朝下执行
	 */
	for {
		//提示客户端设置用户id
		_, err = conn.Write([]byte("请设置用户uid"))
		if err != nil {
			return
		}
		data := make([]byte, 1024)
		c, err := conn.Read(data)
		if err != nil {
			//closed <- true  //这样会阻塞 | 后面取closed的for循环，没有执行到。
			return
		}
		sUid := string(data[0:c])

		//转成int类型
		uid, _ := strconv.Atoi(sUid)
		if uid < 1 {
			continue
		}
		curUid = uid
		talkChan[uid] = make(chan string)
		//fmt.Println(conn, "have set uid ", uid, "can talk")

		_, err = conn.Write([]byte("have set uid " + sUid + " can talk"))
		if err != nil {
			return
		}
		break
	}

	fmt.Println("err 3")

	//当前所有的连接
	fmt.Println(talkChan)

	//读取客户端传过来的数据
	go func() {
		for {
			//不停的读客户端传过来的数据
			data := make([]byte, 1024)
			c, err := conn.Read(data)
			if err != nil {
				fmt.Println("have no client write", err)
				closed <- true //这里可以使用 | 因为是用用的go 新开的线程去处理的。 |  即便chan阻塞，后面的也会执行去读 closed 这个chan
			}

			clientString := string(data[0:c])

			//将客户端过来的数据，写到相应的chan里
			if curUid == 3 {
				talkChan[4] <- clientString
			} else {
				talkChan[3] <- clientString
			}

		}
	}()

	/*
		从chan 里读出给这个客户端的数据 然后写到该客户端里
	*/
	go func() {
		for {
			talkString := <-talkChan[curUid]
			_, err = conn.Write([]byte(talkString))
			if err != nil {
				closed <- true
			}
		}
	}()

	/*
	   检查是否已经关闭连接 如果关闭则推出该线程  去执行defer语句
	*/
	for {
		if <-closed {
			return
		}
	}
}

func main() {

	/**
	建立监听链接
	*/
	ln, err := net.Listen("tcp", "127.0.0.1:1024")
	if err != nil {
		panic(err)
	}

	//创建一个管道

	//talkChan := map[f]
	talkChan := make(map[int]chan string)

	fmt.Printf("%p\n", talkChan)

	/*
	   监听是否有客户端过来的连接请求
	*/
	for {
		fmt.Println("wait connect...")
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("get client connection error: ", err)
		}

		//go handleConnection(conn, talkChan)
		go ClientLogic(conn)
	}
}

func ClientLogic(conn net.Conn) {

	// 从客户端接受数据
	s, _ := bufio.NewReader(conn).ReadString('\n')
	println("由客户端发来的消息：", s)

	// 发送消息给客户端
	conn.Write([]byte("东东你好\n"))

	// 关闭连接
	conn.Close()
	////////cececece
}
