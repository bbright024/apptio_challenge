// When run, sits on a specified port and responds with
// a specified log file.
//
// Usage: logserver -p ":PORT" -d "/directory/of/log/" -f "logfile_name" -addr "ip address for server"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var logdir = flag.String("d", "./", "directory containing app log files")
var port = flag.String("p", ":8888", "log server port")
var addr = flag.String("addr", "localhost", "ip address of server")
var filename = flag.String("f", "mainapp.log", "name of main app log file")

func main() {
	flag.Parse()

	http.HandleFunc("/read", readLog)
	log.Fatal(http.ListenAndServe(*addr+*port, nil))
}

// reads the entire log file and prints it to the socket
func readLog(w http.ResponseWriter, r *http.Request) {
	logfile, err := os.Open(*logdir + *filename)
	if os.IsNotExist(err) {
		fmt.Println(w, "The log file does not exist")
		return
	} else if err != nil {
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
