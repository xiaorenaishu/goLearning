package week_1

import (
	"fmt"
	"testing"
)

/*
	协程
	它比线程更小，十几个goroutine可能体现在底层就是五六个线程，
	Go语言内部帮你实现了这些goroutine之间的内存共享。
	执行goroutine只需极少的栈内存(大概是4~5KB)，当然会根据相应的数据伸缩。
	也正因为如此，可同时运行成千上万个并发任务。goroutine比thread更易用、更高效、更轻便。
	goroutine是通过Go的runtime管理的一个线程管理器。

	默认情况下，channel接收和发送数据都是阻塞的，除非另一端已经准备好，
	这样就使得Goroutines同步变的更加的简单，而不需要显式的lock。
	所谓阻塞，也就是如果读取（value := <-ch）它将会被阻塞，直到有数据接收。
	其次，任何发送（ch<-1）将会被阻塞，直到数据被读出。
*/

func sum(a []int, ch chan int) {
	for _, v := range a {
		fmt.Println(v)
	}
	ch <- 1
}

func TestGoroutine(t *testing.T) {
	var a []int
	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ch := make(chan int)
	fmt.Println("开启协程")
	go sum(a, ch)
	<-ch //读取channel数据，显式阻塞
	fmt.Println("输出处理结果完成")
}

const(
	start int = 1
	end int = 10
)

func ping(pings <-chan int, pongs chan<- int, done chan<- bool, isFirst bool) {
	if isFirst {
		fmt.Printf("I am ping coroutine and start to work\n")
	}

	msg := <-pings
	if msg < end {
		fmt.Printf("I am ping coroutine,receive %d from pong\n", msg)
		pongs <- msg+1
		fmt.Printf("I am ping coroutine,send %d to pong\n", msg+1)
		ping(pings, pongs, done, false)
	} else {
		fmt.Printf("I am ping coroutine,count to %d finish\n", msg)
		done<- true
	}
}

func pong(pings chan<- int, pongs <-chan int, done chan<- bool, isFirst bool) {
	if isFirst {
		fmt.Printf("I am pong coroutine and start to work\n")
	}

	msg := <-pongs
	if msg < end {
		fmt.Printf("I am pong coroutine,receive %d from ping\n", msg)
		pings <- msg+1
		fmt.Printf("I am pong coroutine,send %d to ping\n", msg+1)
		pong(pings, pongs, done, false)
	} else {
		fmt.Printf("I am pong coroutine,count to %d finish\n", msg)
		done<- true
	}
}

func TestGoroutineGroup(t *testing.T) {
	pings := make(chan int, 1)
	pongs := make(chan int, 1)
	done := make(chan bool)

	go ping(pings, pongs, done, true)
	go pong(pings, pongs, done, true)
	fmt.Printf("Coroutine test start %d end %d\n", start, end)
	pings<- start

	<-done
}

