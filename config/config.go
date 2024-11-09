package config

import (
	"os"
)

var Config struct {
	APIURL string
	APIKey string
	Secure bool
}

func GetHost() string {
	if Config.Secure {
		return "https://" + Config.APIURL
	} else {
		return "http://" +  Config.APIURL
	}
}

func GetWsHost() string {
	if Config.Secure {
		return "wss://" + Config.APIURL
	} else {
		return "ws://" + Config.APIURL
	}
}

func GetApiKey()string {
	apk := os.Getenv("HEIMDAHL_API_KEY")
	if apk == "" {
		apk = Config.APIKey
	}
	return apk
}


