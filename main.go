package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)	//connected clients
var broadcast = make(chan Message)		//broadcast channel

//configure the upgrader
var upgrader = websocket.Upgrader{}

type Message struct {
	Email		string `json:"email"`
	Username	string `json:"username"`
	Message		string `json:"message"`
}

func main() {
	fs:= http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	//start listening for incoming chat messages
	go handleMessages()

	//start the server on local host port 8000 and log any errors
	log.Println("https server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}


}


func handleConnections(w http.ResponseWriter, r *http.Request) {
	//Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	//Close the connection when the function returns
	defer ws.Close()

	//Register clients
	clients[ws] = true

	for {
		var msg Message

		//Read in a new message as JSON and map it to the Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		//Grab the next message from the broadcast channel
		msg := <-broadcast
		//Send it to every client that is connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
