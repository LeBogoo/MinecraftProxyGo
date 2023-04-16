package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minecraftproxy/packetUtils"
	"os"
)

type Server struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol int    `json:"protocol"`
}

type Config struct {
	Server                    Server                     `json:"server"`
	ProxyPort                 int                        `json:"proxy_port"`
	StartingDisconnectMessage string                     `json:"starting_disconnect_message"`
	OfflineStatusResponse     packetUtils.StatusResponse `json:"offline_status_response"`
	StartingStatusResponse    packetUtils.StatusResponse `json:"starting_status_response"`
}

func LoadConfig(filename string) Config {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var config Config
	json.Unmarshal(byteValue, &config)

	return config
}
