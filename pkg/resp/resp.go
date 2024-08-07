package resp

import (
	"bufio"
	"bytes"
	"errors"
	"log/slog"
	"net"
	"sync"
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

	mu         sync.Mutex
	listener   net.Listener
	activeConn map[*net.Conn]struct{}
}

func (srv *Server) Close() error {
	srv.shuttingDown.Store(true)

	srv.mu.Lock()
	defer srv.mu.Unlock()

	if err := srv.listener.Close(); err != nil {
		return err
	}
	return srv.closeActiveConnLocked()
}

func (srv *Server) closeActiveConnLocked() error {
	var errs []error

	for conn := range srv.activeConn {
		if err := (*conn).Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (srv *Server) Serve(ln net.Listener) error {
	for {
		conn, err := srv.listener.Accept()
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
	srv.listener = l

	return srv.Serve(l)
}

func (srv *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	srv.mu.Lock()
	if srv.activeConn == nil {
		srv.activeConn = make(map[*net.Conn]struct{})
	}
	srv.activeConn[&conn] = struct{}{}
	srv.mu.Unlock()

	defer func() {
		srv.mu.Lock()
		delete(srv.activeConn, &conn)
		srv.mu.Unlock()
	}()

	var args []string

	scanner := bufio.NewScanner(conn)
	scanner.Split(scanCRLFLines)
	for scanner.Scan() {
		args = append(args, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		slog.Warn("failed to scan connection", "remote", conn.RemoteAddr(), "error", err)
	}

	srv.Handler.ServeRESP(args)
}

var crlf = []byte{'\r', '\n'}

func scanCRLFLines(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, crlf); i >= 0 {
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
