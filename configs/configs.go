// Contains logic for parsing conf files that are formatted in JSON

package configs

import (
	"encoding/json"
	"log"
	"os"
)

type Conf struct {
	Dir     string
	Address string
	Port    string
	Logfile string
}

// reads a conf file in json format 
func ReadConfFile(filename string, c *Conf) {
	confFile, err := os.Open(os.Args[1])
	if err == nil && !os.IsNotExist(err) {
		defer confFile.Close()
		err = json.NewDecoder(confFile).Decode(c)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
	
}
