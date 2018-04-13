// Contains logic for parsing conf files that are formatted in JSON

package configs

import (
	"encoding/json"
	"log"
	"os"
	//	"fmt"
)

// adding a string for time parsing needs to be done... hard to do
// without more information of format
type Conf struct {
	Dir     string
	Address string
	Port    string
	Logfile string
	Timefmt string
}

//func (c Conf) String() string {
//	return fmt.Sprintf(string(c))}

// reads a conf file in json format, saves it in the passed-in struct
func ReadConfFile(filename string, c *Conf) error {
	confFile, err := os.Open(filename)
	if err == nil && !os.IsNotExist(err) {
		defer confFile.Close()
		err = json.NewDecoder(confFile).Decode(c)
		if err != nil {
			log.Printf("Error in readconf - decoder: %v", err)
			return err
		}
	} else if err != nil {
		log.Printf("Error in readconf - open: %v", err)
		return err
	}
	return nil

}
