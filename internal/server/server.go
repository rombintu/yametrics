package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/rombintu/yametrics/internal/config"
	"github.com/rombintu/yametrics/internal/storage"
	"github.com/rombintu/yametrics/lib"
)

type Server struct {
	storage *storage.Storage
	router  *http.ServeMux
	config  config.Config
	Log     *slog.Logger
}

func NewServer(config config.Config) *Server {
	storage, err := storage.NewStorage(config.StorageDriver)
	if err != nil {
		panic(err)
	}
	return &Server{
		config:  config,
		router:  http.NewServeMux(),
		storage: storage,
		Log:     lib.SetupLogger(config.Environment),
	}
}

func (s *Server) Start() {
	s.Log.Info("starting server")
	// Open storage
	s.Log.Debug("opening storage", s.storage.Open())
	s.Log.Info(fmt.Sprintf("server started on http://%s:%d", s.config.Server.Listen, s.config.Server.Port))

	s.ConfigureRouter()
	if err := http.ListenAndServe(
		fmt.Sprintf(
			"%s:%d", s.config.Server.Listen, s.config.Server.Port,
		),
		s.router,
	); err != nil {
		s.Log.Error(err.Error())
		panic(err)
	}
}

func (s *Server) Shutdown() {
	// Close storage
	s.Log.Debug("closing storage", s.storage.Close())
	s.Log.Info("shutdown server")
}

func (s *Server) ConfigureRouter() {
	s.Log.Debug("configure router")
	// s.router.HandleFunc("/", middlewareSetHeaders(rootHandler))

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.Method == http.MethodGet {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			metrics, err := json.Marshal(s.storage.GetAllMetrics())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(metrics)
		} else if r.Method == http.MethodPost {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			if strings.HasPrefix(path, "/update/") {
				parts := strings.Split(strings.TrimPrefix(path, "/update/"), "/")
				if len(parts) == 3 {
					metricType := parts[0]
					metricName := parts[1]
					metricValue := parts[2]

					if err := s.storage.WriteMetric(
						metricType, metricName, metricValue); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					//http.StatusOK
					w.WriteHeader(http.StatusOK)
				} else {
					http.Error(w, "Invalid format", http.StatusNotFound)
				}

			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		}
	})

	// // update counter metrics type
	// s.router.HandleFunc("/update/counter", middlewareSetHeaders(func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	// 	}
	// }))
	s.Log.Debug("configure router is complete")
}
