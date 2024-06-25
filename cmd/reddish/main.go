package main

import (
	"log"
	"net/http"

	"github.com/utilyre/reddish/internal/adapters/mapstorage"
	"github.com/utilyre/reddish/internal/adapters/rpc"
	"github.com/utilyre/reddish/internal/app/service"
)

func main() {
	storageRepo := mapstorage.NewMapStorage()
	storageSVC := service.NewStorageService(storageRepo)
	storageHandler := rpc.NewStorageHandler(storageSVC)

	storageSRV := rpc.NewStorageServer(storageHandler)
	log.Fatal(http.ListenAndServe(":5000", storageSRV))
}
