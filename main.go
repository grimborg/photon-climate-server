package main

import (
	"encoding/json"
	"github.com/grimborg/photon-climate-server/broadcaster"
	"github.com/grimborg/photon-climate-server/photon"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type Config struct {
	DeviceId string `envconfig:"device_id"`
	Token    string
}

func main() {
	var config Config
	if err := envconfig.Process("photon", &config); err != nil {
		log.Fatalln(err)
	}
	c := make(chan photon.Measure)
	go photon.Subscribe(c, config.DeviceId, config.Token)
	bc := broadcaster.New()
	http.Handle("/socket.io/", bc.Server)
	go func() {
		for {
			m := <-c
			s, _ := json.Marshal(m)
			log.Printf("received %+v\n", m)
			bc.Broadcast(string(s))
		}
	}()
	log.Println("Listening at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
