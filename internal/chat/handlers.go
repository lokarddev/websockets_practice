package chat

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

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
	Action      string   `json:"action"`
	Message     string   `json:"message"`
	MessageType string   `json:"message_type"`
	Users       []string `json:"users"`
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

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""
	if err = ws.WriteJSON(response); err != nil {
		log.Println(err)
	}
	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
		}
		log.Println("Error")
	}()
	var payload WsPayload
	for {
		if err := conn.ReadJSON(&payload); err != nil {
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			// list all users and send it back
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.Users = users
			BroadcastToAll(response)
		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			users := getUserList()
			response.Users = users
			BroadcastToAll(response)
		}

		//response.Action = "Got Here"
		//response.Message = fmt.Sprintf("Some message was %s", e.Action)
		//BroadcastToAll(response)
	}
}

func getUserList() []string {
	var users []string
	for _, v := range clients {
		users = append(users, v)
	}
	sort.Strings(users)
	return users
}

func BroadcastToAll(response WsJsonResponse) {
	for client := range clients {
		if err := client.WriteJSON(response); err != nil {
			log.Println(err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}
