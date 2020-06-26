package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	// 同理put函数, os.Open打开文件
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	// 把文件写入响应体
	_, _ = io.Copy(w, f)
}
