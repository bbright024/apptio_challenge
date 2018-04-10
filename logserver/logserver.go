// When run, sits on a specified port and responds with
// a specified log file.
//
// Usage: ./logserver "/path/to/conf/file/conf.json"

package main

import (
	"strings"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
	"apptio/configs"
)

type LogEntry struct {
	Logtime string
	Message string
}

// default configuration settings
var conf = configs.Conf{
	Dir:      "./",
	Address:  "localhost",
	Port:     ":8888",
	Logfile: "mainapp.log",
}

func main() {
	// the log servers log file
	lf, err := os.OpenFile("logserver.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer lf.Close()

	log.SetOutput(lf)
	log.Print("Server Initializing")
	defer log.Print("Server Terminated")
	
	if len(os.Args) > 1 {
		configs.ReadConfFile(os.Args[1], &conf)
	}
	log.Printf("Conf file in use: %v", conf)

	
	logfile, err2 := os.Open(conf.Dir + conf.Logfile)
	if os.IsNotExist(err2) {
		fmt.Fprintln(os.Stderr, "The log file does not exist")
		return
	} else if err2 != nil {
		log.Fatal(err2)
	}
	defer logfile.Close()
	
	http.HandleFunc("/", initRequest)
	err3 := http.ListenAndServe(conf.Address + conf.Port, nil)
	log.Printf("Possible conf file error: %v\n", conf)
	log.Fatal(err3)
}


// converts a logfile into an array of LogEntry structs
func convertLogFile(file *os.File) []LogEntry {
	// not sure if the last arg is taking bytes or slots for string pointers.
	// its filling an array of structs that are themselves arrays of 2 strings,
	// so it SHOULD be slots for string pointers.  will have to make sure later.
	var logs = make([]LogEntry, 0, 200) 
	scanner := bufio.NewScanner(file)

	// the log files are in a predetermined format - let's grab it all
	// and turn each entry into a struct. 
	for scanner.Scan() {
		temp := strings.Split(scanner.Text(), ", ")
		le := LogEntry{Logtime:temp[0], Message:strings.Join(temp[1:], ", ")}
		logs = append(logs, le)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from logfile:", err)
		return nil
	}
	return logs
}

// prints an array of log entries to an io.Writer interface
func printLogs(w io.Writer, logs []LogEntry) {

	fmt.Fprintf(w, " %-9.9s\tMessage\n", "Date")
	
	for _, le := range logs {
		fmt.Fprintf(w, "#%-9.9s\t%s\n", le.Logtime, le.Message) 
	}
}

func initRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connection request from %v: %s", r.RemoteAddr, r.URL.Path)
	logfile, err := os.Open(conf.Dir + conf.Logfile)

	if err != nil {
		log.Fatal(err)
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


// given a date, lists every message in the log file from that date
//func searchLog(w http.ResponseWriter, r * http.Request) {
	
//}


// reads a log file and prints it to the socket buffer
//func readLog(w http.ResponseWriter, r *http.Request, f *os.File) {

//	logs := convertLogFile(f)
//	printLogs(w, logs)
//}
