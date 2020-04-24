package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"net/http"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
	pass   bool // 是否经过密码认证
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func runServer() {
	//flag.Parse()

	debugLog("启动一个 ShareClipServer...")
	debugLog("当前监听端口为:", *port)
	debugLog("当前连接密码为:", *linkKey)
	go manager.start()
	http.HandleFunc("/ws", wsPage)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		debugLog(err.Error())
	}
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			//jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			//manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				//jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				//manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				// 只给通过认证的 conn 发送消息
				if conn.pass {
					select {
					case conn.send <- message:
					default:
						close(conn.send)
						delete(manager.clients, conn)
					}
				}
			}
		}
	}
}

//// server 群发给 server
//func (manager *ClientManager) send(message []byte, ignore *Client) {
//	for conn := range manager.clients {
//		if conn != ignore {
//			conn.send <- message
//		}
//	}
//}

// 不断的从 socket 连接里面读取消息
func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		//debugLog("socket 收到" + string(message))
		var msg SocketMsg
		err = json.Unmarshal(message, &msg)
		if err != nil {
			debugLog("json 解析 socket 消息失败", err)
			debugLog(string(message))
		} else {
			if msg.Key != *linkKey {
				debugLog("连接密码错误, sender:", msg.Sender, ", key:", msg.Key)
			} else {
				// 只将通过认证的消息发送到 broadcast channel 里面
				if msg.Content == MsgContentLinkStart {
					debugLog("认证" + c.id)
					// 如果是认证的话, 将这个 c.pass 设为 true
					c.pass = true
				} else if msg.Content == MsgContentHeartBeat {
					// 心跳包认证
				} else {
					msg.Sender = c.id + ":" + msg.Sender
					jsonMessage, _ := json.Marshal(msg)
					manager.broadcast <- jsonMessage
				}
			}
		}
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()
	// for 循环
	// select c.send 看看有什么需要发送的
	// 如果有需要发送的话就写到这个 socket 里面
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	uid := uuid.NewV1()
	client := &Client{
		id:     uid.String(),
		socket: conn,
		send:   make(chan []byte),
		pass:   false,
	}

	manager.register <- client

	go client.read()
	go client.write()
}
