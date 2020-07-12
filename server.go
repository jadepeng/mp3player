package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Command struct {
	Command string `json:"command"`
	Arg     string `json:"arg"`
}

var mp3Chan = make(chan string, 10)
var quit = make(chan int)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	// cross domain
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)

		var command Command
		decodeErr := json.Unmarshal(message, &command)
		if decodeErr != nil {
			fmt.Println(decodeErr)
			break
		}

		switch command.Command {
		case "play":
			mp3Chan <- command.Arg
			break
		default:
			break
		}

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func play(filename string) {
	log.Printf("start play : %s", filename)

	var err error

	var dec *minimp3.Decoder
	var data, file []byte

	if file, err = ioutil.ReadFile(filename); err != nil {
		log.Fatal(err)
	}

	if dec, data, err = minimp3.DecodeFull(file); err != nil {
		log.Fatal(err)
	}

	c, err := oto.NewContext(dec.SampleRate, dec.Channels, 2, 10240)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	player := c.NewPlayer()
	defer player.Close()

	player.Write(data)

	<-time.After(time.Second)
	dec.Close()
	player.Close()
}

func mp3Player(c chan string, quit chan int) {
	for {
		select {
		case mp3file:= <-c:
			play(mp3file)
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	go mp3Player(mp3Chan,quit)

	http.HandleFunc("/ws", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
