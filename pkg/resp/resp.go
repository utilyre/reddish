package resp

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
)

type Handler interface {
	ServeRESP()
}

type Server struct {
	Addr    string
	Handler Handler

	ln net.Listener
}

func (srv *Server) Close() error {
	// TODO: close open connections
	return srv.ln.Close()
}

func (srv *Server) Serve(ln net.Listener) error {
	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			slog.Warn("failed to accept connection", "error", err)
			continue
		}

		go srv.handleConn(conn)
	}
}

func (srv *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	srv.ln = l

	return srv.Serve(l)
}

func (srv *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Printf("'%s'\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		slog.Warn("failed to scan connection", "remote", conn.RemoteAddr(), "error", err)
	}
}
