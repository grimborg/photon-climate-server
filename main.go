package main

import (
	"github.com/grimborg/photon-climate-server/photon"
	"github.com/kelseyhightower/envconfig"
	"log"
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
	var c chan photon.Measure
	go photon.Subscribe(c, config.DeviceId, config.Token)
	for {
		m := <-c
		log.Printf("received %+v\n", m)
	}
}
