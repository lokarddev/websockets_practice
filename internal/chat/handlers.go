package chat

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var views = jet.NewSet(jet.NewOSFileSystemLoader("./templates"), jet.InDevelopmentMode())
var upgradeConn = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Home(w http.ResponseWriter, r *http.Request) {
	if err := renderPage(w, "home.html", nil); err != nil {
		log.Println(err)
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	if err = view.Execute(w, data, nil); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines response back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndpoint upgrades connection to a endpoint
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConn.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to endpoint.")
	var response WsJsonResponse
	response.Message = `<em><small>Connected to a server.</small></em>`

	if err = ws.WriteJSON(response); err != nil {
		log.Println(err)
	}
}
