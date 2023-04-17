package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minecraftproxy/minecraft"
	"minecraftproxy/packetUtils"
	"os"
)

type Server struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol int    `json:"protocol"`
}

type StatusResponses struct {
	Offline  packetUtils.StatusResponse `json:"offline"`
	Starting packetUtils.StatusResponse `json:"starting"`
}

type DisconnectMessages struct {
	NowStarting minecraft.Text `json:"now_starting"`
	Starting    minecraft.Text `json:"starting"`
}

type Config struct {
	Server             Server             `json:"server"`
	ApproxStartupTime  int                `json:"approx_startup_time"`
	ProxyPort          int                `json:"proxy_port"`
	StatusResponses    StatusResponses    `json:"status_responses"`
	DisconnectMessages DisconnectMessages `json:"disconnect_messages"`
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
