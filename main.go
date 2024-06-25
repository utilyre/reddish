package main

import (
	"context"
	"log"
	"net/http"

	"github.com/utilyre/reddish/rpc/reddish"
)

func main() {
	storage := NewMapStorage()
	tsrv := reddish.NewReddishServer(NewReddishService(storage))

	log.Fatal(http.ListenAndServe(":5000", tsrv))
}

type ReddishService struct {
	storage Storage
}

func NewReddishService(storage Storage) *ReddishService {
	return &ReddishService{storage: storage}
}

func (rs *ReddishService) Set(ctx context.Context, r *reddish.SetReq) (*reddish.SetResp, error) {
	err := rs.storage.Set(ctx, r.Key, r.Val)
	if err != nil {
		return nil, err
	}

	return &reddish.SetResp{}, nil
}

func (rs *ReddishService) Get(ctx context.Context, r *reddish.GetReq) (*reddish.GetResp, error) {
	val, err := rs.storage.Get(ctx, r.Key)
	if err != nil {
		return nil, err
	}

	return &reddish.GetResp{Val: val}, nil
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
