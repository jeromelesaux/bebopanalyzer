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
	confFile, errOpen := os.Open(file)

	if errOpen != nil {
		fmt.Println("Error while opening file ", file, errOpen.Error())
	}
	defer confFile.Close()
	decoder := json.NewDecoder(confFile)
	errDecode := decoder.Decode(conf)
	if errDecode != nil {
		fmt.Println("Cannot decode configuration file." + errDecode.Error())
	}
	return conf
}
