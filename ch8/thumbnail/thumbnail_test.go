// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// This file is just a place to put example code from the book.
// It does not actually run any code in gopl.io/ch8/thumbnail.

package thumbnail_test

import (
	"log"
	"os"
	"sync"

	"gopl.io/ch8/thumbnail"
)

//!+1
// makeThumbnails makes thumbnails of the specified files.
func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

//!-1

//!+2
// NOTE: incorrect!
func makeThumbnails2(filenames []string) {
	for _, f := range filenames {
		go thumbnail.ImageFile(f) // NOTE: ignoring errors
	}
}

//!-2

//!+3
// makeThumbnails3 makes thumbnails of the specified files in parallel.
// makeThumbnails3 并行生成指定文件的缩略图
func makeThumbnails3(filenames []string) {
	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f) // NOTE: ignoring errors
			ch <- struct{}{}
		}(f)
	}

	// Wait for goroutines to complete.
	for range filenames {
		<-ch
	}
}

//!-3

//!+4
// makeThumbnails4 makes thumbnails for the specified files in parallel.
// It returns an error if any step failed.
func makeThumbnails4(filenames []string) error {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := thumbnail.ImageFile(f)
			errors <- err
		}(f)
	}

	for range filenames {
		if err := <-errors; err != nil {
			return err // NOTE: incorrect: goroutine leak!
		}
	}

	return nil
}

//!-4

//!+5
// makeThumbnails5 makes thumbnails for the specified files in parallel.
// It returns the generated file names in an arbitrary order,
// or an error if any step failed.
func makeThumbnails5(filenames []string) (thumbfiles []string, err error) {
	type item struct {
		thumbfile string
		err       error
	}

	// 使用一个缓冲通道来返回生成的图像文件的名称以及任何错误消息。
	ch := make(chan item, len(filenames))

	for _, f := range filenames {
		go func(f string) {
			var it item
			// 接收缩略图文件和错误
			it.thumbfile, it.err = thumbnail.ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return thumbfiles, nil
}

//!-5

//!+6
// makeThumbnails6 makes thumbnails for each file received from the channel.
// It returns the number of bytes occupied by the files it creates.
func makeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // number of working goroutines
	for f := range filenames {
		// wg.Add(1) 必须在工作 goroutine 开始之前执行，而不是在中间
		wg.Add(1)
		// worker
		go func(f string) {
			// wg.Done() 等价于 wg.Add(-1)
			// 使用 defer 确保发生错误的情况下计数器可以递减
			defer wg.Done()
			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb) // OK to ignore error
			sizes <- info.Size()
		}(f)
	}

	// closer
	/*
		在关闭 sizes 通道之前，等待所有worker结束。
		这里等待和关闭必须和在sizes通道上面的迭代并行执行。
		如果将等待操作放在循环之前的主 goroutine 中，因为通道会满，它不会结束；
		如果放在循环后面，它将不可达，因为没有任何东西可用来关闭通道，循环可能永不结束。
	*/
	go func() {
		wg.Wait()
		close(sizes)
	}()

	// chan sizes 将每个文件的大小带回主 goroutine
	// 使用 range 循环进行接收然后计算总和
	// 主 goroutine 把大多数时间花在了 range 循环休眠上，等待工作者发送或等待 closer 来关闭通道。
	var total int64
	for size := range sizes {
		total += size
	}
	return total
}

//!-6
