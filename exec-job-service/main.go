package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Ror testing only
var (
	statusByJobID = map[string]string{}
	mu            sync.RWMutex
)

func main() {
	http.HandleFunc("/run", handleRun)
	http.HandleFunc("/status/{jobId}", handleStatus)
	http.ListenAndServe(":4000", nil)
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	jobID, err := createID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeAsJSON(w, map[string]any{
			"msg":  err.Error(),
			"code": http.StatusInternalServerError,
		})
		return
	}
	setJobStatus(jobID, "started")

	go func() {
		// Simulate running a command

		id := 1
		for id < 30 {
			setJobStatus(jobID, fmt.Sprintf("running ID: %d", id))
			time.Sleep(1 * time.Second)
			id += 1
		}
		setJobStatus(jobID, fmt.Sprintf("completed ID: %d", id))
	}()

	writeAsJSON(w, map[string]any{
		"jobId": jobID,
	})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("jobId")
	writeAsJSON(w, map[string]any{
		"jobId":  jobID,
		"status": getJobStatus(jobID),
	})
}

func createID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b[:24]); err != nil {
		return "", fmt.Errorf("creating ID: %w", err)
	}

	// 8 bytes used for the timestamp i.e. 32 - 8 = 24
	binary.BigEndian.PutUint64(b[24:], uint64(time.Now().UnixMilli()))
	return hex.EncodeToString(b), nil
}

func getJobStatus(jobID string) string {
	mu.RLock()
	defer mu.RUnlock()

	status, ok := statusByJobID[jobID]
	if !ok {
		return "not-found"
	}
	return status
}

func setJobStatus(jobID, status string) {
	mu.Lock()
	defer mu.Unlock()

	statusByJobID[jobID] = status
}

func writeAsJSON(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")

	// Ignore the error
	json.NewEncoder(w).Encode(res)
}
