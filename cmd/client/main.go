package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/utilyre/reddish/internal/adapters/rpc"
)

func main() {
	ctx := context.Background()

	client := rpc.NewStorageProtobufClient("http://localhost:6979", http.DefaultClient)

	_, err := client.Set(ctx, &rpc.SetReq{Key: "some", Val: []byte("something")})
	if err != nil {
		log.Fatal(err)
	}

	delResp, err := client.Del(ctx, &rpc.DelReq{Keys: []string{"some"}})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("deleted", delResp.NumDeleted, "keys")

	getResp, err := client.Get(ctx, &rpc.GetReq{Key: "some"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(getResp.Val))
}
