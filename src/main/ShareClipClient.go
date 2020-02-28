package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gorilla/websocket"
	"net/url"
	"strconv"
	"time"
)

var addr = flag.String("addr", "localhost:7878", "http service address")
var debug = flag.Bool("debug", false, "print debug log")
var linkKeyClient = flag.String("key", "ShareClip", "the link key about ShareClip Server")
var continueLink = flag.Bool("continue", false, "Constantly try to reconnect in disconnected, otherwise it will only try a limited times")
var clientVersion = flag.Bool("v", false, "show the version")

const ClientVerString string = "ShareClip Client  V1.0.1"

var localClipTemp = ""

type SocketMsg struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func main() {
	flag.Parse()

	if *clientVersion {
		fmt.Println(ClientVerString)
		return
	}

	// 启动一个 Web Socket
	fmt.Println("连接:", *addr)
	fmt.Println("连接密码为:", *linkKeyClient)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("连接成功")
	}

	go clipListen(conn)
	go sendHeartbeat(conn)

	var sockeMsgJson SocketMsg
	errorCount := 0
	maxErrorCount := 120
	errorWaitTime := 30
	// 读取信息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log("socket 读取错误 :" + err.Error())
			//return
			// 如果不是不断的尝试的话， 就计算重连次数
			if !*continueLink {
				errorCount += 1
				if errorCount >= maxErrorCount {
					fmt.Println("连接断开")
					return
				}
			}
			log(strconv.Itoa(errorWaitTime) + "s 后尝试重新连接")
			time.Sleep(time.Second * 30)
		}
		_ = json.Unmarshal(message, &sockeMsgJson)
		if sockeMsgJson.Content != localClipTemp {
			log("接收到 clip 更新 : " + sockeMsgJson.Content)
			log("更新本地剪贴板")
			localClipTemp = sockeMsgJson.Content
			_ = clipboard.WriteAll(sockeMsgJson.Content)
		}
	}
}

// 监听剪贴板
func clipListen(conn *websocket.Conn) {
	var clipTemp string
	var err error
	for {
		clipTemp, err = clipboard.ReadAll()
		if err == nil && clipTemp != "" && clipTemp != localClipTemp {
			localClipTemp = clipTemp
			log("更新 localClipTemp 为 : " + clipTemp)
			//conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
			_ = conn.WriteMessage(websocket.TextMessage, []byte(clipTemp))
			log("写入 Websocket 里")
		}

		time.Sleep(time.Millisecond * 1200)
	}
}

func sendHeartbeat(conn *websocket.Conn) {
	for {
		time.Sleep(time.Minute * 5)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(""))
		log("发送一个心跳包")
	}
}

func log(msg string) {
	if *debug {
		fmt.Println(msg)
	}
}
