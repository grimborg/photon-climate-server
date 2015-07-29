package main

import (
	"encoding/json"
	"fmt"
	"github.com/grimborg/photon-climate-server/broadcaster"
	"github.com/grimborg/photon-climate-server/photon"
	"github.com/grimborg/photon-climate-server/storage"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type Config struct {
	DeviceId  string `envconfig:"device_id"`
	Token     string
	RedisPort int    `envconfig:"redis_port"`
	RedisHost string `envconfig:"redis_host"`
}

func main() {
	var config Config
	if err := envconfig.Process("photon", &config); err != nil {
		log.Fatalln(err)
	}
	c := make(chan photon.Measure)
	go photon.Subscribe(c, config.DeviceId, config.Token)
	archive := storage.New(config.RedisHost, config.RedisPort)
	getHistory := func(w http.ResponseWriter, r *http.Request) {
		history, err := archive.ReadAll()
		if err != nil {
			log.Fatal("Error getting the history", err)
		}
		data, err := json.Marshal(history)
		log.Println(data)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprint(w, string(data))
	}
	bc := broadcaster.New()
	http.Handle("/socket.io/", bc.Server)
	http.HandleFunc("/history/", getHistory)
	go func() {
		for {
			m := <-c
			s, _ := json.Marshal(m)
			log.Printf("received %+v\n", m)
			archive.Add(m)
			bc.Broadcast(string(s))
		}
	}()
	log.Println("Listening at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
