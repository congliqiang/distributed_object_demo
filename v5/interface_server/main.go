package main

import (
	"github.com/congliqiang/distributed_object_demo/v5/interface_server/heartbeat"
	"github.com/congliqiang/distributed_object_demo/v5/interface_server/locate"
	"github.com/congliqiang/distributed_object_demo/v5/interface_server/objects"
	"github.com/congliqiang/distributed_object_demo/v5/interface_server/versions"
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
