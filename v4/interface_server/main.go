package main

import (
	"distributed_object_demo/v4/interface_server/heartbeat"
	"distributed_object_demo/v4/interface_server/locate"
	"distributed_object_demo/v4/interface_server/objects"
	"distributed_object_demo/v4/interface_server/versions"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
