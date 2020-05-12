// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 219.
//!+

// Clock1 is a TCP server that periodically writes the time.

/* 
	顺序时钟服务器，它以每秒钟一次的频率向客户端发送当前时间
*/
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		// 一次性只能处理一个客户请求，
		// 第二个客户端必须等到第一个结束才能正常工作，
		// 因为服务器是顺序的
		handleConn(conn) // handle one connection at a time
	}
}

// 处理一个完整的客户连接
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

//!-
