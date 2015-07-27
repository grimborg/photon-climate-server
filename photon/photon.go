package photon

import (
	"encoding/json"
	"github.com/donovanhide/eventsource"
	"log"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Data        string    `json:"data"`
	PublishedAt time.Time `json:"published_at"`
}

type Measure struct {
	Timestamp   time.Time `json:"timestamp"`
	Temperature int64     `json:"temperature"`
	Humidity    int64     `json:"humidity"`
}

func (m Message) Measure() Measure {
	parts := strings.Split(m.Data, ":")
	temperature, _ := strconv.ParseInt(parts[0], 10, 64)
	humidity, _ := strconv.ParseInt(parts[0], 10, 64)
	return Measure{
		Temperature: temperature,
		Humidity:    humidity,
		Timestamp:   m.PublishedAt,
	}
}

func Subscribe(c chan Measure, deviceId string, token string) {
	url := "https://api.particle.io/v1/devices/" +
		deviceId +
		"/events/temperature?access_token=" +
		token
	stream, err := eventsource.Subscribe(url, "")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		event := <-stream.Events
		if event.Event() != "temperature" {
			continue
		}
		message := Message{}
		if err := json.Unmarshal([]byte(event.Data()), &message); err != nil {
			log.Fatalln(err)
		}
		measure := message.Measure()
		log.Printf("%+v\n", measure)
		c <- measure
	}
}
