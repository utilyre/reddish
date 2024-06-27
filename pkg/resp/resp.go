package resp

import (
	"bufio"
	"bytes"
	"errors"
	"log/slog"
	"net"
	"sync/atomic"
)

var (
	ErrServerClosed = errors.New("server closed")
)

type Handler interface {
	ServeRESP(args []string)
}

type HandlerFunc func(args []string)

func (f HandlerFunc) ServeRESP(args []string) {
	f(args)
}

type Server struct {
	Addr    string
	Handler Handler

	shuttingDown atomic.Bool
	ln           net.Listener
}

func (srv *Server) Close() error {
	srv.shuttingDown.Store(true)
	// TODO: close open connections
	return srv.ln.Close()
}

func (srv *Server) Serve(ln net.Listener) error {
	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			if srv.shuttingDown.Load() {
				return ErrServerClosed
			}

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

	var args []string

	scanner := bufio.NewScanner(conn)
	scanner.Split(ScanCRLFLines)
	for scanner.Scan() {
		args = append(args, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		slog.Warn("failed to scan connection", "remote", conn.RemoteAddr(), "error", err)
	}

	srv.Handler.ServeRESP(args)
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
