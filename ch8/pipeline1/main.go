// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 228.

// Pipeline1 demonstrates an infinite 3-stage pipeline.
package main

import "fmt"

//!+
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; ; x++ {
			naturals <- x
		}
	}()

	// Squarer
	go func() {
		// 死循环
		for {
			x := <-naturals
			squares <- x * x
		}
	}()

	// 主函数的死循环
	// Printer (in main goroutine)
	for {
		fmt.Println(<-squares)
	}
}

//!-
