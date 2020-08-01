package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// PUT /objects/<object_name>HTTP/1.1
func put(w http.ResponseWriter, r *http.Request) {
	// r.URL.EscapedPath() 返回经过转义的后的路径部分 形式为/objects/<object_name>, os.Create创建文件
	f, e := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	// 将对象体存入文件f
	_, _ = io.Copy(f, r.Body)
}
