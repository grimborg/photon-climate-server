package storage

import (
	"encoding/json"
	"github.com/grimborg/photon-climate-server/photon"
	"gopkg.in/redis.v3"
	"strconv"
)

type storage struct {
	redis *redis.Client
}

func New(host string, port int) *storage {
	h := &storage{
		redis: redis.NewClient(&redis.Options{
			Addr:     host + ":" + strconv.Itoa(port),
			Password: "",
			DB:       0,
		}),
	}
	return h
}

func (s storage) Add(m photon.Measure) error {
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = s.redis.LPush("measures", string(j)).Result()
	return err
}

func (s storage) ReadAll() ([]photon.Measure, error) {
	data, err := s.redis.LRange("measures", 0, 10).Result()
	if err != nil {
		return nil, err
	}
	measures := make([]photon.Measure, len(data))
	for idx, d := range data {
		m := photon.Measure{}
		err := json.Unmarshal([]byte(d), &m)
		if err != nil {
			return nil, err
		}
		measures[idx] = m
	}
	return measures, nil
}
