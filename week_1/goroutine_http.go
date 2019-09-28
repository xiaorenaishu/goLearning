package main

import (
	"fmt"
	"net"
	"os"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"sync"
	"bufio"
	"runtime"
	"bytes"
	"strconv"
)

// 日志列表最大长度
const maxLogLength int = 3
// 日志列表
var logList = make([]string, 0)
// 日志刷写互斥锁
var logFlushLock = &sync.RWMutex{}

func main() {
	address := ":8080"

	// 定时器
	cronTicker := time.NewTicker(time.Millisecond * 200)
	go func() {
		for range cronTicker.C {
			// 定时刷写日志
			flushLog()
		}
	}()

	// http.HandleFunc("/", handleHTTPRequest)
	// print("Listening from tcp " + address)
	// http.ListenAndServe(address, nil)

	// 开始网络监听
	print("Listening from tcp " + address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		print(err.Error())
		return
	}

	// 监听kill信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		print("Receive exit signal " + sig.String())
		flushLog()
		cronTicker.Stop()
		// todo 平滑关闭，确保所有请求都处理完成后才关闭
		listener.Close()
	}()

	// 开始接收请求
	for {
		client, err := listener.Accept()

		if err != nil {
			print(err.Error())
			break
		}

		go handleRequest(client)
	}
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf(
		"Receive http request from remote addr:%s",
		r.RemoteAddr)

	log(msg)
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, time.Now().Format("2006-01-02 15:04:05.0000") + " Hello World")
}

func handleRequest (client net.Conn) {
	if client == nil {
		return
	}

	requestByte := make([]byte, 1024)
	log(fmt.Sprintf("Receive connection from remote addr:%s", client.RemoteAddr()))
	_, err := client.Read(requestByte[:])
	if err != nil {
		print(err.Error())
		return
	}

	log(fmt.Sprintf("Receive data:\n%s", string(requestByte)))

	// 模拟响应耗时
	time.Sleep(10 * time.Second)
	client.Write([]byte("HTTP/1.1 200 OK\n\n" + time.Now().Format("2006-01-02 15:04:05.0000") + " Hello World"))
	client.Close()
}

/**
 * 记录日志
 */
func log(log string) {
	log = fmt.Sprintf("[%d:%d] [%s] %s", os.Getpid(), getGid(), time.Now().Format("2006-01-02 15:04:05.0000"), log)

	logFlushLock.RLock()
	logList = append(logList, log)
	logFlushLock.RUnlock()

	// 日志缓冲区饱和则主动执行刷写
	if len(logList) >= maxLogLength {
		go flushLog()
	}
}

/**
 * 输出日志
 */
func print(log string) {
	log = fmt.Sprintf("[%d:%d] [%s] %s", os.Getpid(), getGid(), time.Now().Format("2006-01-02 15:04:05.0000"), log)
	fmt.Println(log)
}

/**
 * 刷写日志
 */
func flushLog() {
	if (len(logList) == 0) {
		return
	}

	// 使用互斥锁防止日志列表
	logFlushLock.Lock()
	logListCopy := logList
	logList = make([]string, 0)
	logFlushLock.Unlock()

	f, err := os.OpenFile("/Users/cxr/go/src/goLearning/week_1/week_2.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		print(err.Error())
	}
	// 程序执行后关闭文件
	defer f.Close()

	// 使用带缓冲区的写入方式，减少磁盘I/O
	w := bufio.NewWriter(f)
	for _, logString := range logListCopy {
		// 无日志则停止
		if logString == "" {
			break
		}

		// 模拟日志写入阻塞
		time.Sleep(1 * time.Second)
		if _, err := w.WriteString(logString + "\n"); err != nil {
			print(err.Error())
		}
	}

	print("flush log")
	// 确保所有数据从操作系统缓冲区中写入磁盘
	w.Flush()
}

func getGid() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
