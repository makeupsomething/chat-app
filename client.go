package main

import (
	"time"

	"github.com/gorilla/websocket"
)

//client is a single chatting user
type client struct {
	//socket is the web socket for this user
	socket *websocket.Conn

	//send is the channel on which the messages are sent
	send chan *message

	//room is the "chatroom" the user is in
	room *room

	//userData holds information about the user
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarUrl, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarUrl.(string)
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
