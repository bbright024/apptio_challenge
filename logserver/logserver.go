// When run, sits on a specified port and responds with
// a specified log file.
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

// made global for testing purposes - helps keep tests from brittleness
var msgdatefmt = " Date\t\t\tMessage\n"

func main() {
	// the log servers log file
	lf, err := os.OpenFile("./logserver.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
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

	logfile, err2 := os.Open(conf.Dir + conf.Logfile)
	if os.IsNotExist(err2) {
		fmt.Fprintln(os.Stderr, "The log file does not exist")
		return
	} else if err2 != nil {
		log.Fatal(err2)
	}
	defer logfile.Close()

	http.HandleFunc("/", initRequest)
	err3 := http.ListenAndServe(conf.Address+conf.Port, nil)
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
		// if the datetime format of the logtime section were known,
		// I could change the LogEntry struct to be
		// {Logtime time.Time, Message string}
		// but life isn't fair.

		// presently, this assumes that the logtime section of the line
		// has no commas.  if it does, this will break!!
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

// prints an array of log entries to an io.Writer interface
func printLogs(w io.Writer, logs []LogEntry) {

	fmt.Fprintf(w, msgdatefmt)

	for _, le := range logs {
		fmt.Fprintf(w, "#%-9.9s\t%s\n", le.Logtime, le.Message)
	}
}

// there's a lot of options here.  every call to the log server needs to at least
// print to the logserver log file with the connection request, so to me it seems that
// every request should go through a dispatcher function that matches requests with
// a switch statement jump table.  however, "TGPL" specifically calls this out as bad
// style.  What's a coder to do?  I don't want redundant code, yet I don't want a
// switch statement with a million cases.
// i guess what bugs me is that the http.handler stuff doesn't seem much different than
// a switch statement.  not sure what the benefit is, because i can't see what's under the
// hood in that http.Handle/HandleFunc call.
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

// given a date, lists every message in the log file from that date
//func searchLog(w http.ResponseWriter, r * http.Request) {

//}

// reads a log file and prints it to the socket buffer
//func readLog(w http.ResponseWriter, r *http.Request, f *os.File) {

//	logs := convertLogFile(f)
//	printLogs(w, logs)
//}
