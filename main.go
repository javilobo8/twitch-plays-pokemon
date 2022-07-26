package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/micmonay/keybd_event"
)

type Action struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type Config struct {
	Channel  string   `json:"channel"`
	Interval int      `json:"interval"`
	Keys     []Action `json:"keys"`
}

var (
	client               *twitch.Client
	config               *Config
	lastCommandTimestamp = time.Now().UnixMilli()
)

func loadConfig() {
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &config)
}

func findAction(message string) *Action {
	for _, item := range config.Keys {
		if item.Message == message {
			return &item
		}
	}
	return nil
}

func pressKey(action *Action) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	event := 0x00
	switch action.Action {
	case "BUTTON_UP":
		event = keybd_event.VK_D
	case "BUTTON_DOWN":
		event = keybd_event.VK_F
	case "BUTTON_LEFT":
		event = keybd_event.VK_G
	case "BUTTON_RIGHT":
		event = keybd_event.VK_H
	case "BUTTON_START":
		event = keybd_event.VK_ENTER
	case "BUTTON_SELECT":
		event = keybd_event.VK_BACKSPACE
	case "BUTTON_A":
		event = keybd_event.VK_Z
	case "BUTTON_B":
		event = keybd_event.VK_X
	case "BUTTON_L":
		event = keybd_event.VK_A
	case "BUTTON_R":
		event = keybd_event.VK_S
	}

	if event != 0x00 {
		kb.Clear()
		kb.SetKeys(event)
		kb.Press()
		time.Sleep(100 * time.Millisecond)
		kb.Release()
		kb.Clear()
	} else {
		log.Println("ERROR receiving event")
	}
}

func OnPrivateMessage(message twitch.PrivateMessage) {
	if strings.ToLower(config.Channel) == strings.ToLower(message.Channel) {
		now := time.Now().UnixMilli()
		if lastCommandTimestamp+int64(config.Interval*1000) < now {
			lastCommandTimestamp = now
			action := findAction(message.Message)

			if action != nil {
				log.Println(fmt.Sprintf("%s => %s", message.User.Name, action.Action))
				pressKey(action)
			}
		}
	}
}

func initServices() {
	// Twitch
	client = twitch.NewAnonymousClient()
	client.OnPrivateMessage(OnPrivateMessage)
	client.OnConnect(func() {
		fmt.Println("[Twitch] Connected")
	})
	fmt.Println(fmt.Sprintf("[Twitch] JOIN %s", config.Channel))
	client.Join(config.Channel)
	fmt.Println("[Twitch] Connecting")
	client.Connect()
}

func main() {
	loadConfig()
	initServices()
}
