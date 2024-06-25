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

	client := rpc.NewStorageProtobufClient("http://localhost:5500", http.DefaultClient)

	_, err := client.Set(ctx, &rpc.SetReq{Key: "some", Val: []byte("something")})
	if err != nil {
		log.Fatal(err)
	}

	// _, err := client.Delete(ctx, &rpc.DeleteReq{Key: "some"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := client.Get(ctx, &rpc.GetReq{Key: "some"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp.Val))
}
