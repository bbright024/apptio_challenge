// When run, sits on a specified port and responds with
// a specified log file.
//
// Usage: logserver "/path/to/conf/file/conf.json"

package main

import (
	"strings"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
)

type Conf struct {
	Dir     string
	Address string
	Port    string
	Logfile string
}

type LogEntry struct {
	Logtime string
	Message string
}

// default configuration settings
var conf = Conf{
	Dir:      "./",
	Address:  "localhost",
	Port:     ":8888",
	Logfile: "mainapp.log",
}

const (
	MaxLogEntries = 200
)

var logfile *os.File

func main() {
	
	if len(os.Args) > 1 {
		confFile, err := os.Open(os.Args[1])
		if err == nil && !os.IsNotExist(err) {
			defer confFile.Close()

			err = json.NewDecoder(confFile).Decode(&conf)

			if err != nil {
				log.Fatal(err)
			}
			
		} else if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Conf file in use: ")
	fmt.Print(conf)
	logfile, err2 := os.Open(conf.Dir + conf.Logfile)
	if os.IsNotExist(err2) {
		fmt.Fprintln(os.Stderr, "The log file does not exist")
		return
	} else if err2 != nil {
		log.Fatal(err2)
	}
	defer logfile.Close()
	
	http.HandleFunc("/read", readLog)
	err3 := http.ListenAndServe(conf.Address + conf.Port, nil)
	fmt.Fprintf(os.Stderr, "Possible conf file error: %v\n", conf)
	log.Fatal(err3)
}

// reads the entire log file and prints it to the socket buffer
func readLog(w http.ResponseWriter, r *http.Request) {
	logfile, err := os.Open(conf.Dir + conf.Logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()

	// not sure if the last arg is taking bytes or slots for string pointers.
	// its filling an array of structs that are themselves arrays of 2 strings,
	// so it SHOULD be slots for string pointers.  will have to make sure later.
	var logs = make([]LogEntry, 0, MaxLogEntries) 
	scanner := bufio.NewScanner(logfile)

	// the log files are in a predetermined format - let's grab it all
	// and turn each entry into a struct. 
	for scanner.Scan() {
		temp := strings.Split(scanner.Text(), ", ")
		le := LogEntry{Logtime:temp[0], Message:strings.Join(temp[1:], ", ")}
		logs = append(logs, le)
	}

	fmt.Fprintf(w, " %-9.9s\tMessage\n", "Date")
	for _, le := range logs {
		fmt.Fprintf(w, "#%-9.9s\t%s\n", le.Logtime, le.Message) 
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from logfile:", err)
	}
}
