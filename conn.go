package ftp

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

// Conn represents a connection to the FTP server
type Conn struct {
	server   *Server
	rwc      net.Conn
	dataType dataType
	dataPort *dataPort
	workDir  string
	rootDir  string
}

func (c *Conn) serve() {
	defer c.rwc.Close()
	log.Println("<<", c.rwc.RemoteAddr().String(), "connected!")
	// send welcome
	c.respond(status220)
	// read line from connection
	scanner := bufio.NewScanner(c.rwc)
	for scanner.Scan() {
		input := strings.Fields(scanner.Text())
		if len(input) == 0 {
			continue
		}
		cmd, args := input[0], input[1:]
		log.Println(cmd, args)

		switch cmd {
		case "PWD":
			c.pwd(args)
		case "CWD": // cd
			c.cd(args)
		case "LIST": // ls
			c.ls(args)
		case "PORT":
			c.port(args)
		case "USER":
			c.user(args)
		case "QUIT": // close
			c.respond(status221)
			return
		case "RETR": // get
			c.retr(args)
		case "LPRT":
			c.setDataType(args)
		default:
			c.respond(status502)
		}
	}
}

func (c *Conn) dataConnect() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.dataPort.toAddress())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Conn) respond(s string) {
	log.Println(">>", s)
	_, err := io.WriteString(c.rwc, s+c.EOL())
	if err != nil {
		log.Println(err)
	}
}

// EOL returns the line terminator matching the FTP standard for the datatype.
func (c *Conn) EOL() string {
	switch c.dataType {
	case ascii:
		return "\r\n"
	case binary:
		return "\n"
	default:
		return "\n"
	}
}
