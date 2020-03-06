package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	registry   = make(map[string]*Item)
	registryMu = &sync.RWMutex{}
)

type Item struct {
	Port string
	Seen int64
}

func getURL(ch, file string) *string {
	registryMu.Lock()
	defer registryMu.Unlock()

	item, ok := registry[ch]

	if !ok {
		return nil
	}

	if item.Seen < time.Now().Unix() {
		return nil
	}

	url := fmt.Sprintf("http://localhost:%s/%s", item.Port, file)
	return &url
}

func ping(ch, port string) {
	registryMu.Lock()
	defer registryMu.Unlock()
	_, ok := registry[ch]
	if !ok {
		registry[ch] = &Item{
			Port: port,
			Seen: time.Now().Unix() + 3,
		}
	}
	registry[ch].Port = port
	registry[ch].Seen = time.Now().Unix() + 3

}
