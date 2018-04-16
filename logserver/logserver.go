// When run, listens on a port specified in the passed-in config file
// and responds to GET requests with a log file also specified in the config file.
//
// Usage: ./logserver "/path/to/conf/file/conf.json"

package main

import (
	"apptio/configs"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// TODO: Seperate data structure type defs from main
type LogEntry struct {
	Logtime string
	Message string
}

// default configuration settings
var conf = configs.Conf{
	Dir:     "./",
	Address: "localhost",
	Port:    ":8888",
	Logfile: "mainapp.log",
	Timefmt: "",
}

func main() {
	// For security, the logserver's log file must be created prior to execution
	lf, err := os.OpenFile("./logserver.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Please create a log file for the logserver before executing.")
		} 
		log.Fatal(err)
	}
	defer lf.Close()

	log.SetOutput(lf)
	log.Print("Server Initializing")
	defer log.Print("Server Terminated")

	if len(os.Args) > 1 {
		err = configs.ReadConfFile(os.Args[1], &conf)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Conf file in use: %v", conf)

	// Check to see if the main app log file is available
	mainAppLogFile, err2 := os.Open(conf.Dir + conf.Logfile)
	if err2 != nil {
		if os.IsNotExist(err2) {
			log.Printf("The main app log file does not exist")
		}
		log.Fatal(err2)
	}
	defer mainAppLogFile.Close()

	// TODO: If caching is an option, read and parse the main app log file here to a global
	//       to keep from having to read the file at every /read request from clients
	
	http.HandleFunc("/", initRequest)
	err3 := http.ListenAndServe(conf.Address+conf.Port, nil)
	log.Printf("Possible conf file error: %v\n", conf)
	log.Fatal(err3)
}

// Converts a logfile into an array of LogEntry structs
func convertLogFile(file *os.File) []LogEntry {
	// alloc some space for the logentry array - prevents early resizing in append()
	var logs = make([]LogEntry, 0, 200)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		//                     ****** WARNING *******
		//  This code assumes the logtime section of the line has no commas.
		//  If it does, output might differ from what is expected.
		temp := strings.Split(scanner.Text(), ", ")
		le := LogEntry{Logtime: temp[0], Message: strings.Join(temp[1:], ", ")}
		logs = append(logs, le)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from logfile:", err)
		return nil
	}
	return logs
}

// Format string for log table header 
var msgdatefmt = " Date\t\t\tMessage\n"

// Prints an array of log entries to an io.Writer interface
func printLogs(w io.Writer, logs []LogEntry) {
	fmt.Fprintf(w, msgdatefmt)
	for _, le := range logs {
		fmt.Fprintf(w, "#%-9.9s\t%s\n", le.Logtime, le.Message)
	}
}

// The handler for our server:
//   initializes common data and uses a switch statement to direct requests
func initRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connection request from %v: %s", r.RemoteAddr, r.URL.Path)
	logfile, err := os.Open(conf.Dir + conf.Logfile)

	if err != nil {
		log.Printf("Error: %v", err)
	}
	defer logfile.Close()

	// turn the log file into an array of log entries
	var logs []LogEntry
	logs = convertLogFile(logfile)
	if logs == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Logfile conversion error, sorry\n")
		return
	}

	switch r.URL.Path {
	case "/read":
		w.WriteHeader(http.StatusAccepted)
		printLogs(w, logs)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such page: %s\n", r.URL)
	}
}

