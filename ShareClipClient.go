package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

var addr = flag.String("addr", "localhost:7878", "http service address")
var debug = flag.Bool("debug", false, "print debug log")

var localClipTemp = ""

type SocketMsg struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func main() {
	flag.Parse()
	// 启动一个 Web Socket
	fmt.Println("连接 : " + *addr)
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

	var sockeMsgJson SocketMsg
	// 读取信息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
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
	for {
		clipTemp, _ = clipboard.ReadAll()
		if clipTemp != localClipTemp {
			localClipTemp = clipTemp
			log("更新 localClipTemp 为 : " + clipTemp)
			//conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
			_ = conn.WriteMessage(websocket.TextMessage, []byte(clipTemp))
			log("写入 Websocket 里")
		}
		time.Sleep(time.Millisecond * 1200)
	}
}

func log(msg string) {
	if *debug {
		fmt.Println(msg)
	}
}
