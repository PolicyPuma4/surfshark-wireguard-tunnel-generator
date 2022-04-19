// Created by https://github.com/PolicyPuma4
// Repository https://github.com/PolicyPuma4/surfshark-wireguard-tunnel-generator

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Server struct {
	ConnectionName string
	PubKey         string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("With Surfshark connect to a server using the WireGuard protocol")
	fmt.Print("Press ENTER to continue")
	fmt.Scanln()

	data, err := os.ReadFile("C:\\ProgramData\\Surfshark\\WireguardConfigs\\SurfsharkWireGuard.conf")
	check(err)

	lines := strings.Split(strings.ReplaceAll(string(data), "%", "%%"), "\r\n")
	for index, line := range lines {
		if strings.HasPrefix(line, "PublicKey = ") {
			lines[index] = "PublicKey = %s"
			continue
		}
		if strings.HasPrefix(line, "Endpoint = ") {
			u, err := url.Parse("//" + strings.TrimPrefix(line, "Endpoint = "))
			check(err)
			lines[index] = "Endpoint = %s:" + u.Port()
			break
		}
	}

	template := strings.Join(lines, "\r\n")

	servers := []Server{}
	for _, endpoint := range [...]string{"generic", "double", "static", "obfuscated"} {
		resp, err := http.Get("https://api.surfshark.com/v4/server/clusters/" + endpoint + "?countryCode=")
		check(err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		check(err)

		respServers := []Server{}
		err = json.Unmarshal(body, &respServers)
		check(err)

		servers = append(servers, respServers...)
	}

	path := "surfshark-wireguard-tunnel-generator"
	err = os.MkdirAll(path, os.ModeDir)
	check(err)

	processed := 0
	for _, server := range servers {
		if len(server.PubKey) == 0 {
			continue
		}

		f, err := os.Create(path + "\\" + server.ConnectionName + ".conf")
		check(err)
		defer f.Close()

		_, err = f.WriteString(fmt.Sprintf(template, server.PubKey, server.ConnectionName))
		check(err)

		processed++
	}

	fmt.Println(fmt.Sprintf("Created %d files", processed))
	fmt.Print("Press ENTER to exit")
	fmt.Scanln()
}
