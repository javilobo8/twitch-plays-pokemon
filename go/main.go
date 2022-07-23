package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	client      *twitch.Client
	config      *Config
	lastCommand = time.Now().UnixMilli()
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

func OnPrivateMessage(message twitch.PrivateMessage) {
	if config.Channel == message.Channel {
		now := time.Now().UnixMilli()
		if lastCommand+int64(config.Interval*1000) < now {
			lastCommand = now
			log.Println(fmt.Sprintf("%s: %s", message.User.Name, message.Message))
			action := findAction(message.Message)
			if action != nil {
				log.Println(fmt.Sprintf("%s: %s", action.Action, action.Message))
				kb, err := keybd_event.NewKeyBonding()
				if err != nil {
					panic(err)
				}
				switch action.Action {
				case "BUTTON_UP":
					kb.SetKeys(keybd_event.VK_UP)
				case "BUTTON_DOWN":
					kb.SetKeys(keybd_event.VK_DOWN)
				case "BUTTON_LEFT":
					kb.SetKeys(keybd_event.VK_LEFT)
				case "BUTTON_RIGHT":
					kb.SetKeys(keybd_event.VK_RIGHT)
				case "BUTTON_START":
					kb.SetKeys(keybd_event.VK_ENTER)
				case "BUTTON_SELECT":
					kb.SetKeys(keybd_event.VK_BACKSPACE)
				case "BUTTON_A":
					kb.SetKeys(keybd_event.VK_Z)
				case "BUTTON_B":
					kb.SetKeys(keybd_event.VK_X)
				case "BUTTON_L":
					kb.SetKeys(keybd_event.VK_A)
				case "BUTTON_R":
					kb.SetKeys(keybd_event.VK_S)
				}
				kb.Press()
				time.Sleep(10 * time.Millisecond)
				kb.Release()
			}
		}
	}
}

func initServices() {
	// Twitch
	client = twitch.NewAnonymousClient()
	client.OnPrivateMessage(OnPrivateMessage)
	client.OnConnect(func() {
		fmt.Println("Connected")
	})
	fmt.Println(fmt.Sprintf("JOIN %s", config.Channel))
	client.Join(config.Channel)
	fmt.Println("Connecting")
	client.Connect()
}

func main() {
	loadConfig()
	initServices()
}
