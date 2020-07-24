package ftp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (c *Conn) user(args []string) {
	c.respond(fmt.Sprintf(status230, strings.Join(args, " ")))
}

func (c *Conn) pwd(args []string) {
	if len(args) > 0 {
		c.respond(status501)
		return
	}
	c.respond(c.workDir)
}

func (c *Conn) cd(args []string) {
	if len(args) != 1 {
		c.respond(status501)
		return
	}
	workDir := filepath.Join(c.workDir, args[0])
	absPath := filepath.Join(c.rootDir, workDir)
	_, err := os.Stat(absPath)
	if err != nil {
		log.Print(err)
		c.respond(status550)
		return
	}
	c.workDir = workDir
	c.respond(status200)
}

func (c *Conn) ls(args []string) {
	var target string
	if len(args) > 0 {
		target = filepath.Join(c.rootDir, c.workDir, args[0])
	} else {
		target = filepath.Join(c.rootDir, c.workDir)
	}

	files, err := ioutil.ReadDir(target)
	if err != nil {
		log.Print(err)
		c.respond(status550)
		return
	}
	c.respond(status150)

	dataConn, err := c.dataConnect()
	if err != nil {
		log.Print(err)
		c.respond(status425)
		return
	}
	defer dataConn.Close()

	for _, file := range files {
		_, err := fmt.Fprint(dataConn, file.Name(), file.Size(), c.EOL())
		if err != nil {
			log.Print(err)
			c.respond(status426)
		}
	}

	c.respond(status226)
}

func (c *Conn) port(args []string) {
	if len(args) != 1 {
		c.respond(status501)
		return
	}
	dataPort, err := dataPortFromHostPort(args[0])
	if err != nil {
		log.Print(err)
		c.respond(status501)
		return
	}
	c.dataPort = dataPort
	c.respond(status200)
}

func (c *Conn) retr(args []string) {
	if len(args) != 1 {
		c.respond(status501)
		return
	}

	path := filepath.Join(c.rootDir, c.workDir, args[0])
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		c.respond(status550)
	}
	c.respond(status150)

	dataConn, err := c.dataConnect()
	if err != nil {
		log.Print(err)
		c.respond(status425)
	}
	defer dataConn.Close()

	_, err = io.Copy(dataConn, file)
	if err != nil {
		log.Print(err)
		c.respond(status426)
		return
	}
	io.WriteString(dataConn, c.EOL())
	c.respond(status226)
}
