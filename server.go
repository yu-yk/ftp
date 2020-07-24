package ftp

import (
	"log"
	"net"
)

type Server struct {
	Addr    string
	RootDir string
}

func ListenAndServe(addr, rootDir string) error {
	server := &Server{Addr: addr, RootDir: rootDir}
	return server.ListenAndServe()
}

func (srv *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	log.Println("Listening on " + srv.Addr)
	return srv.Serve(listener)
}

func (srv *Server) Serve(listener net.Listener) error {
	for {
		rwConn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		defer rwConn.Close()
		ftpConn := srv.newConn(rwConn)
		go ftpConn.serve()
		// go handleConn(conn)
	}
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) *Conn {
	c := &Conn{
		server:  srv,
		rwc:     rwc,
		workDir: "/",
		rootDir: srv.RootDir,
	}
	return c
}
