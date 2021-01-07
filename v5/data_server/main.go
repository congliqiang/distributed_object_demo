package main

import (
	"github.com/congliqiang/distributed_object_demo/v5/data_server/heartbeat"
	"github.com/congliqiang/distributed_object_demo/v5/data_server/locate"
	"github.com/congliqiang/distributed_object_demo/v5/data_server/objects"
	"github.com/congliqiang/distributed_object_demo/v5/data_server/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
