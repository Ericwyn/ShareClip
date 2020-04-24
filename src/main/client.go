package main

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gorilla/websocket"
	"net/url"
	"strconv"
	"time"
)

var localClipTemp = ""

const MsgContentHeartBeat string = "ShareClip_Msg_HeartBeat"
const MsgContentLinkStart string = "ShareClip_Msg_LinkStart"

func runClient() {
	//flag.Parse()
	// 给这个发送端一个id, 默认使用 unix time
	if *senderName == "" {
		*senderName = fmt.Sprint(time.Now().Unix())
	}

	// 启动一个 Web Socket
	fmt.Println("连接:", *addr)
	fmt.Println("连接密码为:", *linkKey)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("连接成功")
	}

	// 发送一个消息进行认证
	msg := SocketMsg{
		Sender:  *senderName,
		Content: MsgContentLinkStart,
		Key:     *linkKey, // 连接密码
	}
	jsonByte, _ := json.Marshal(msg)
	_ = conn.WriteMessage(websocket.TextMessage, jsonByte)

	go clipListen(conn)
	go sendHeartbeat(conn)

	var socketMsgJson SocketMsg
	errorCount := 0
	maxErrorCount := 120
	errorWaitTime := 30
	// 读取信息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			debugLog("socket 读取错误 :" + err.Error())
			//return
			// 如果不是不断的尝试的话， 就计算重连次数
			if !*continueLink {
				errorCount += 1
				if errorCount >= maxErrorCount {
					fmt.Println("连接断开")
					return
				}
			}
			debugLog(strconv.Itoa(errorWaitTime) + "s 后尝试重新连接")
			time.Sleep(time.Second * 30)
		}
		_ = json.Unmarshal(message, &socketMsgJson)
		if socketMsgJson.Content != localClipTemp {
			debugLog("接收到 clip 更新 : " + socketMsgJson.Content)
			debugLog("更新本地剪贴板")
			localClipTemp = socketMsgJson.Content
			_ = clipboard.WriteAll(socketMsgJson.Content)
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
			debugLog("更新 localClipTemp 为 : " + clipTemp)
			//conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
			msg := SocketMsg{
				Sender:  *senderName,
				Content: clipTemp,
				Key:     *linkKey, // 连接密码
			}
			jsonByte, _ := json.Marshal(msg)
			// json 格式化之后发出
			_ = conn.WriteMessage(websocket.TextMessage, jsonByte)
			debugLog("写入 Websocket 里")
		}

		time.Sleep(time.Millisecond * 1200)
	}
}

func sendHeartbeat(conn *websocket.Conn) {
	for {
		time.Sleep(time.Minute * 5)
		msg := SocketMsg{
			Sender:  *senderName,
			Content: MsgContentHeartBeat,
			Key:     *linkKey, // 连接密码
		}
		jsonByte, _ := json.Marshal(msg)
		_ = conn.WriteMessage(websocket.TextMessage, jsonByte)
		debugLog("发送一个心跳包")
	}
}

func debugLog(msg ...interface{}) {
	if *debug {
		fmt.Println(msg)
	}
}
