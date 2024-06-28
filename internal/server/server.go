package server

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rombintu/yametrics/internal/config"
	"github.com/rombintu/yametrics/internal/storage"
	"github.com/rombintu/yametrics/lib"
)

type Server struct {
	storage *storage.Storage
	router  *echo.Echo
	config  config.ServerConfig
	Log     *slog.Logger
}

func NewServer(config config.ServerConfig) *Server {
	storage, err := storage.NewStorage(config.StorageDriver)
	if err != nil {
		panic(err)
	}
	return &Server{
		config: config,
		// Change router. Increment3
		router:  echo.New(),
		storage: storage,
		Log:     lib.SetupLogger(config.Environment),
	}
}

func (s *Server) Start() {
	s.Log.Info("starting server", slog.String("address", s.config.Address))
	// Open storage
	s.storage.Open()
	s.Log.Debug("opening storage")

	s.ConfigureRouter()

	if err := http.ListenAndServe(
		s.config.Address,
		s.router,
	); err != nil {
		s.Log.Error(err.Error())
		panic(err)
	}
	// s.router.Start(fmt.Sprintf(
	// 	"%s:%d", s.config.Server.Listen, s.config.Server.Port,
	// ))
}

func (s *Server) Shutdown() {
	// Close storage
	s.Log.Debug("closing storage", s.storage.Close())
	s.Log.Info("shutdown server")
}

func (s *Server) ConfigureRouter() {
	s.Log.Debug("configure router")
	// s.router.HandleFunc("/", middlewareSetHeaders(rootHandler))
	// s.router.GET("/", rootHandler)
	s.router.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, s.storage.GetAllMetrics())
	})

	s.router.GET("/value/:mtype/:mname", s.getMetricsByNames)

	s.router.POST("/update/:mtype/:mname/:mvalue", s.updateMetrics)

	s.Log.Debug("configure router is complete")
}
