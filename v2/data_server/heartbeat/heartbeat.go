package heartbeat

import (
	"distributed_object_demo/v2/rabbitmq"
	"os"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBIT_MQ_SERVER"))
	defer q.Close()
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
