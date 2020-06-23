package main

import (
	"distributed_object_demo/v2/interface_server/heartbeat"
	"distributed_object_demo/v2/interface_server/locate"
	"distributed_object_demo/v2/interface_server/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
