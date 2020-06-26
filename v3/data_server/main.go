package main

import (
	"distributed_object_demo/v3/data_server/heartbeat"
	"distributed_object_demo/v3/data_server/locate"
	"distributed_object_demo/v3/data_server/objects"
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
