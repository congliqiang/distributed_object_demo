package main

import (
	"github.com/congliqiang/distributed_object_demo/v3/data_server/heartbeat"
	"github.com/congliqiang/distributed_object_demo/v3/data_server/locate"
	"github.com/congliqiang/distributed_object_demo/v3/data_server/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
