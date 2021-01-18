package rs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/congliqiang/distributed_object_demo/V6/utils"
	"github.com/congliqiang/distributed_object_demo/v6/interface_server/objectstream"
	"github.com/klauspost/reedsolomon"
	"io"
	"log"
	"net/http"
)

const (
	DATA_SHARDS     = 4
	PARITY_SHARDS   = 2
	ALL_SHARDS      = DATA_SHARDS + PARITY_SHARDS
	BLOCK_PER_SHARD = 8000
	BLOCK_SIZE      = BLOCK_PER_SHARD * DATA_SHARDS
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	for i := range writers {
		writers[i], e = objectstream.NewTempPutStream(dataServers[i], fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	enc := NewEncoder(writers)
	return &RSPutStream{enc}, nil
}

func (s *RSPutStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0
	for length != 0 {
		next := BLOCK_SIZE - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	shards, _ := e.enc.Split(e.cache)
	e.enc.Encode(shards)
	for i := range shards {
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}

type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, e := NewRSPutStream(dataServers, hash, size)
	if e != nil {
		return nil, e
	}
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	token := &resumableToken{name, size, hash, dataServers, uuids}
	return &RSResumablePutStream{putStream, token}, nil
}

func (s *RSResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(b)
}

func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}
	var t resumableToken
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}

	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &objectstream.TempPutStream{t.Servers[i], t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumablePutStream{&RSPutStream{enc}, &t}, nil
}

func (s *RSResumablePutStream) CurrentSize() int64 {
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(r.Header) * DATA_SHARDS
	if size > s.Size {
		size = s.Size
	}
	return size
}
