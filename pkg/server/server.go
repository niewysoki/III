package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rzetelskik/iii/pkg/store"
	"k8s.io/klog/v2"
	"net/http"
)

type stat struct {
	Seen     bool `json:"seen"`
	Clicked  bool `json:"clicked"`
	LoggedIn bool `json:"logged_in"`
}

var (
	st store.Store[string, stat] = store.NewThreadSafeMapStore[string, stat]()
)

func userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if !values.Has("u") {
			http.Error(w, "required parameter 'u' is missing", http.StatusBadRequest)
			return
		}

		user := values.Get("u")
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
	})

}

func indexMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		klog.V(4).InfoS("received webpage request", "User", user)

		stats, ok := st.Get(user)
		if !ok {
			stats = stat{}
		}
		stats.Clicked = true

		st.Put(user, stats)

		next.ServeHTTP(w, r)
	})
}

func pixelMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		klog.V(4).InfoS("received pixel request", "User", user)

		stats, ok := st.Get(user)
		if !ok {
			stats = stat{}
		}
		stats.Seen = true

		st.Put(user, stats)

		next.ServeHTTP(w, r)
	})
}

func loginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := r.FormValue("u")
		if len(user) < 1 {
			http.Error(w, "missing 'u' in form", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		if len(username) < 1 {
			http.Error(w, "missing 'username' in form", http.StatusBadRequest)
			return
		}

		password := r.FormValue("password")
		if len(password) < 1 {
			http.Error(w, "missing 'password' in form", http.StatusBadRequest)
			return
		}

		klog.V(4).InfoS("received login form", "User", user, "Username", username, "Password", password)

		stats, ok := st.Get(user)
		if !ok {
			stats = stat{}
		}
		stats.LoggedIn = true

		st.Put(user, stats)

		next.ServeHTTP(w, r)
	})
}

func handleLogin(redirectAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectAddr, http.StatusSeeOther)
	}
}

func handleStats(w http.ResponseWriter, _ *http.Request) {
	stats := st.GetAll()

	payload, err := json.Marshal(stats)
	if err != nil {
		err = fmt.Errorf("can't marshal response: %w", err)
		klog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func NewServer(addr, redirectAddr string) *http.Server {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./assets/server/static"))

	r.Handle("/pixel.gif", userMiddleware(pixelMiddleware(fs))).
		Methods(http.MethodGet)

	r.HandleFunc("/stats", handleStats).
		Methods(http.MethodGet)

	r.Handle("/", userMiddleware(indexMiddleware(fs))).
		Methods(http.MethodGet)

	r.Handle("/", loginMiddleware(handleLogin(redirectAddr))).
		Methods(http.MethodPost)

	r.PathPrefix("/").Handler(fs)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
