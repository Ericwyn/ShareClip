package main

type SocketMsg struct {
	Sender  string `json:"sender,omitempty"`
	Content string `json:"content,omitempty"`
	Key     string `json:"key"`
}
