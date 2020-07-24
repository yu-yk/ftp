package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yu-yk/ftp"
)

var port int
var rootDir string

func init() {
	flag.IntVar(&port, "p", 8080, "port number")
	flag.StringVar(&rootDir, "d", "public", "root directory")
	flag.Parse()
}

func main() {
	err := ftp.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), rootDir)
	if err != nil {
		log.Fatalln(err)
	}
}
