package main

import (
	"log"
	"net/http"

	"github.com/utilyre/reddish/internal/adapters/mapstorage"
	"github.com/utilyre/reddish/internal/adapters/rpc"
	"github.com/utilyre/reddish/internal/app/service"
	"github.com/utilyre/reddish/rpc/storage"
)

func main() {
	storageRepo := mapstorage.NewMapStorage()
	storageSVC := service.NewStorageService(storageRepo)
	storageHandler := rpc.NewStorageHandler(storageSVC)

	tsrv := storage.NewStorageServer(storageHandler)
	log.Fatal(http.ListenAndServe(":5000", tsrv))
}

/*
type Server struct {
	Addr   string
	quitCh chan struct{}
	ln     net.Listener
}

func NewServer(addr string) *Server {
	return &Server{Addr: addr}
}

func (srv *Server) Start() error {
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}
	srv.ln = ln

	return srv.acceptConns()
}

func (srv *Server) Close() error {
	close(srv.quitCh)
	return nil
}

func (srv *Server) acceptConns() error {
	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			slog.Info("could not accept connection", "remote", conn.RemoteAddr(), "error", err)
			continue
		}

		go srv.handleConn(conn)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
	defer func() {
		slog.Info("closing connection", "remote", conn.RemoteAddr())
		if err := conn.Close(); err != nil {
			slog.Warn("failed to close connection", "remote", conn.RemoteAddr(), "error", err)
		}
	}()

	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanWords)

	fmt.Println()
	for scanner.Scan() {
		fmt.Printf("'%s'\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		slog.Warn("failed to read connection", "error", err)
		return
	}
}
*/
