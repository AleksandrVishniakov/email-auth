package configs

import (
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/repositories/postgres"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/servers"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

const configPath = "./configs/config.yml"

type Configs struct {
	HTTP *servers.HTTPConfigs `yaml:"http"`
	DB   *postgres.DBConfigs  `yaml:"db"`
}

func MustConfigs() *Configs {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var bytes []byte
	var buffer = make([]byte, 256)

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			bytes = append(bytes, buffer[:n]...)
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		bytes = append(bytes, buffer[:n]...)
	}

	var configs Configs

	err = yaml.Unmarshal(bytes, &configs)
	if err != nil {
		log.Fatal(err)
	}

	return &configs
}
