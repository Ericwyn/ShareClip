package main

import (
	"flag"
	"fmt"
)

// server flag
var port = flag.String("port", "7878", "[server] 服务监听端口")

// client flag
var addr = flag.String("addr", "localhost:7878", "[client] 连接的地址")
var continueLink = flag.Bool("continue", false, "[client] 与服务断开连接之后是否持续重连, 否则只重连有限次数")
var senderName = flag.String("sender", "", "[client] 客户端标记名")

// 通用 flag
var runServerFlag = flag.Bool("server", false, "运行 ShareClip Server")
var runClientFlag = flag.Bool("client", false, "运行 ShareServer Client")

var debug = flag.Bool("debug", false, "打印 debug debugLog")
var linkKey = flag.String("key", "share", "连接密码")

var version = flag.Bool("v", false, "版本号")

const versionCode = "v2.1"

func main() {
	flag.Parse()
	if *version {
		fmt.Println("ShareClip", versionCode)
		return
	}
	if *runServerFlag {
		runServer()
	} else if *runClientFlag {
		runClient()
	} else {
		fmt.Println("请选择运行 Server 或 Client 服务")
	}
}
