// Contains logic for parsing conf files that are formatted in JSON

package configs

import (
	"encoding/json"
	"log"
	"os"
)

// adding a string for time parsing needs to be done... hard to do
// without more information of format
type Conf struct {
	Dir      string
	Address  string
	Port     string
	Logfile  string
	Timefmt string
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
