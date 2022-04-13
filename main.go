// Created by https://github.com/PolicyPuma4
// Repository https://github.com/PolicyPuma4/surfshark-wireguard-tunnel-generator

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	ConnectionName string
	PubKey         string
}

func main() {
	defer fmt.Scanln()

	configPath := "C:\\ProgramData\\Surfshark\\WireguardConfigs\\SurfsharkWireGuard.conf"
	fmt.Println("Open Surfshark and connect to a server using the WireGuard protocol.")
	fmt.Print("Press ENTER to continue...")
	fmt.Scanln()
	if _, err := os.Stat(configPath); err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	configContents := buf.String()

	endpoints := [...]string{"generic", "double", "static", "obfuscated"}
	servers := []Server{}
	for _, endpoint := range endpoints {
		resp, err := http.Get("https://api.surfshark.com/v4/server/clusters/" + endpoint + "?countryCode=")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		respServers := []Server{}
		err = json.Unmarshal(body, &respServers)
		if err != nil {
			log.Fatal(err)
		}

		servers = append(servers, respServers...)
	}

	privateKey := ""
	for _, line := range strings.Split(configContents, "\r\n") {
		if strings.HasPrefix(line, "PrivateKey = ") {
			privateKey = strings.TrimPrefix(line, "PrivateKey = ")
			break
		}
	}
	if len(privateKey) == 0 {
		log.Fatal("Unable to find private key.")
	}

	outputDir := "Surfshark WireGuard"
	if _, err := os.Stat(outputDir); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(outputDir, os.ModeDir)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	for _, server := range servers {
		if len(server.PubKey) == 0 {
			continue
		}

		f, err := os.Create(outputDir + "\\" + server.ConnectionName + ".conf")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		s := strings.ReplaceAll(fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = 10.14.0.2/16
DNS = 162.252.172.57, 149.154.159.92
[Peer]
PublicKey = %s
AllowedIps = 0.0.0.0/0
Endpoint = %s:51820
[Peer]
PublicKey = o07k/2dsaQkLLSR0dCI/FUd3FLik/F/HBBcOGUkNQGo=
AllowedIPs = 172.16.0.36/32
Endpoint = 92.249.38.1:51820
`, privateKey, server.PubKey, server.ConnectionName), "\n", "\r\n")

		f.WriteString(s)
	}

	fmt.Println("Complete.")
	fmt.Print("Press ENTER to exit...")
}
