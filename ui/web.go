package ui

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"drtelemetry/telemetry"
	"github.com/gobuffalo/packr"
	"fmt"
)

var Addr = flag.String("http", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var connections []*websocket.Conn

var box = packr.NewBox("../template")
var ui = box.String("ui.html")
var homeTemplate = template.Must(template.New("").Parse(ui))

func ListenAndServe(dataChannel chan telemetry.TelemetryData) {
	http.HandleFunc("/ws", HandleWs)
	http.HandleFunc("/", HandleUi)

	go func() {
		for {
			data := <-dataChannel
			for _, foo := range connections {
				foo.WriteJSON(data)
			}
		}
	}()
	fmt.Printf("Listening for clients on on %s\n", *Addr)
	http.ListenAndServe(*Addr, nil)
}

func HandleWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	connections = append(connections, c)
	select {}
}

func HandleUi(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}