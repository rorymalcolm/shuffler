package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var logger *zap.Logger

const workerPoolSize = 10

type Worker struct {
	ID int
}

func (worker *Worker) NewWorker(id int) *Worker {
	return &Worker{ID: id}
}

func (worker *Worker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "worker with id=%v\n", worker.ID)
}

func main() {
	logger = createLogger()
	defer logger.Sync()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /path/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "got path\n")
	})

	mux.HandleFunc("/task/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "handling task with id=%v\n", id)
	})

	workpool := make(map[int]*Worker)
	for i := 0; i < workerPoolSize; i++ {
		workpool[i] = &Worker{ID: i}
	}

	// simulate a failure by removing a worker
	delete(workpool, 3)

	mux.HandleFunc("/worker/", func(w http.ResponseWriter, r *http.Request) {
		shuffledShards := deriveShuffledShards(workerPoolSize, r)
		for _, shard := range shuffledShards {
			worker, ok := workpool[shard]
			if !ok {
				logger.Error("worker not found", zap.Int("shard", shard))
				http.Error(w, "worker not found", http.StatusInternalServerError)
				return
			}
			worker.ServeHTTP(w, r)
		}
	})

	http.ListenAndServe("localhost:8090", mux)
}

func deriveShuffledShards(shards int, r *http.Request) []int {
	logger.Info("deriving shuffled shards")
	ip := r.RemoteAddr
	logger.Info("ip", zap.String("ip", ip))
	hash := md5.Sum([]byte(ip))
	hashInt := binary.BigEndian.Uint64(hash[:])
	shard := int(hashInt % uint64(shards))
	return []int{shard}
}
