// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"os"
	"time"
)

//!+

/*
	使用 select 多路复用
*/
func main() {
	// ...create abort channel...

	//!-

	//!+abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	//!-abort

	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")
	// select 像 switch 语句一样，它有一系列的情况和一个可选的默认分支。
	// 每一个情况指定一次通信（在一些通道上进行发送或接收操作），
	// 或者在一个短变量声明中
	select {
	case <-time.After(10 * time.Second):	// 指示事件过去10s
		// Do nothing.
		fmt.Println("After 10s!")
	case <-abort:	// 中止事件
		fmt.Println("Launch aborted!")
		return
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
