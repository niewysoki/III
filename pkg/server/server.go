package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	serverassets "github.com/rzetelskik/iii/assets/server"
	"github.com/rzetelskik/iii/pkg/store"
	"k8s.io/klog/v2"
	"net/http"
)

type stat struct {
	Seen    bool `json:"seen"`
	Clicked bool `json:"clicked"`
}

var (
	tsms *store.ThreadSafeMapStore[string, stat]
)

func init() {
	tsms = store.NewThreadSafeMapStore[string, stat]()
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if !values.Has("u") {
		klog.V(4).Infof("received webpage request without 'u' parameter")
		http.Error(w, "required parameter 'u' is missing", http.StatusBadRequest)
		return
	}

	user := values.Get("u")
	klog.V(4).Infof("received webpage request for user %s", user)

	// TODO
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	stats, ok := tsms.Get(user)
	if !ok {
		stats = stat{}
	}
	stats.Clicked = true

	tsms.Put(user, stats)
}

func handlePixel(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if !values.Has("u") {
		klog.V(4).Infof("received pixel request without 'u' parameter")
		http.Error(w, "required parameter 'u' is missing", http.StatusBadRequest)
		return
	}

	user := values.Get("u")
	klog.V(4).Infof("received pixel request for user %s", user)

	w.Header().Add("Content-Type", "image/gif")

	_, err := w.Write(serverassets.Pixel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stats, ok := tsms.Get(user)
	if !ok {
		stats = stat{}
	}
	stats.Seen = true

	tsms.Put(user, stats)
}

func handleStats(w http.ResponseWriter, _ *http.Request) {
	stats := tsms.GetAll()
	m, ok := stats.(map[string]stat)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(m)
	if err != nil {
		err = fmt.Errorf("can't marshal response: %w", err)
		klog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func NewServer(addr string) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/", handlePage).
		Methods(http.MethodGet)

	r.HandleFunc("/pixel", handlePixel).
		Methods(http.MethodGet)

	r.HandleFunc("/stats", handleStats).
		Methods(http.MethodGet)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
