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
	fs:= http.FileServer(http.Dir("../public"))
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
