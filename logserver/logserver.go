package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
//	"net/url"
	"os"
//	"strings"
)

var logdir = flag.String("d", "/home/bbright/", "directory containing app log files")
var port = flag.String("p", ":8888", "log server port")
var addr = flag.String("addr", "localhost", "ip address of server")
var filename = flag.String("f", "mainapp.log", "name of main app log file")

func main() {
	flag.Parse()



	http.HandleFunc("/readLog", readLog)
	log.Fatal(http.ListenAndServe(*addr+*port, nil))
}

func readLog(w http.ResponseWriter, r *http.Request) {
	logfile, err := os.Open(*logdir + *filename)
	if err != nil {
		log.Fatal(err)
	}

	defer logfile.Close()
	scanner := bufio.NewScanner(logfile)
	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading from logfile:", err)
	}
}

//func search(w http.ResponseWriter, r *http.Request) {

//}

// opens a log file and returns a map of dates to
//func parseLog() {

//}
