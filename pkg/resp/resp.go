package resp

import (
	"bufio"
	"bytes"
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
	scanner.Split(ScanCRLFLines)
	for scanner.Scan() {
		fmt.Printf("'%s'\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		slog.Warn("failed to scan connection", "remote", conn.RemoteAddr(), "error", err)
	}
}

func ScanCRLFLines(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
