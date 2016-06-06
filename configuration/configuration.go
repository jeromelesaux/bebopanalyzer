package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type AppConfiguration struct {
	HttpPort        string
	BasepathStorage string
}

func LoadConfiguration(file string) *AppConfiguration {
	conf := new(AppConfiguration)
	confFile, err := os.Open(file)

	if err != nil {
		fmt.Println("Error while opening file ", file, err.Error())
	}
	defer confFile.Close()
	decoder := json.NewDecoder(confFile)
	err = decoder.Decode(conf)
	if err != nil {
		fmt.Println("Cannot decode configuration file." + err.Error())
	}
	return conf
}
