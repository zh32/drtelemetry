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

var box = packr.NewBox("../resources")
var ui = box.String("ui.html")
var homeTemplate = template.Must(template.New("").Parse(ui))

func ListenAndServe(dataChannel chan telemetry.TelemetryData) {
	http.HandleFunc("/ws", HandleWs)
	http.Handle("/assets/", http.FileServer(box))
	http.Handle("/custom/", http.StripPrefix("/custom/", http.FileServer(http.Dir("."))))
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

type TemplateData struct {
	CustomStyle   string
	WebsocketHost string
}

func HandleUi(w http.ResponseWriter, r *http.Request) {
	templateData := TemplateData{r.URL.Query().Get("style"), "ws://" + r.Host + "/ws"}

	homeTemplate.Execute(w, templateData)
}
