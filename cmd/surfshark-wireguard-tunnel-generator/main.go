package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

type Server struct {
	Name string `json:"connectionName"`
	Key  string `json:"pubKey"`
}

func main() {
	prompt := promptui.Prompt{
		Label: "Private key",
		Validate: func(input string) (err error) {
			if len(input) == 0 {
				return errors.New("no private key provided")
			}

			_, err = base64.StdEncoding.DecodeString(input)
			return
		},
		Mask:        '*',
		HideEntered: true,
	}

	privateKey, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	template := `[Interface]
PrivateKey = %s
Address = 10.14.0.2/16
DNS = 162.252.172.57, 149.154.159.92

[Peer]
PublicKey = %s
AllowedIPs = 0.0.0.0/0
Endpoint = %s:51820
`

	prompt = promptui.Prompt{
		Label: "Output directory",
		Default: func() (dir string) {
			dir, _ = os.Getwd()
			return
		}(),
		AllowEdit: true,
		Validate: func(input string) (err error) {
			_, err = os.Stat(input)
			return
		},
	}

	outputDirectory, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range [...]string{"generic", "double", "static", "obfuscated"} {
		resp, err := http.Get("https://api.surfshark.com/v4/server/clusters/" + path)
		if err != nil {
			log.Fatal(err)
		}
		defer func(resp *http.Response) {
			resp.Body.Close()
		}(resp)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		connections := []Server{}
		err = json.Unmarshal(body, &connections)
		if err != nil {
			log.Fatal(err)
		}

		for _, connection := range connections {
			if len(connection.Key) == 0 {
				continue
			}

			filePrefix, _ := strings.CutSuffix(connection.Name, ".prod.surfshark.com")
			file, err := os.Create(filepath.Join(outputDirectory, filePrefix+".conf"))
			if err != nil {
				log.Fatal(err)
			}
			defer func(file *os.File) {
				file.Close()
			}(file)

			_, err = file.WriteString(fmt.Sprintf(template, privateKey, connection.Key, connection.Name))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
