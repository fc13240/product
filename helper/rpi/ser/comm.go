package ser

import (
	"helper/configs"
	"net/http"

	"github.com/gorilla/websocket"
)

type Receive struct {
	link   *websocket.Conn
	Cmmand string    `json:"cmmand"`
	Args   configs.M `json:"args"`
}

func (receive *Receive) Say(content string) {
	receive.link.WriteJSON(configs.M{"cmmand": "text", "args": configs.M{"content": content}})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
