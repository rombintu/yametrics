package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/rombintu/yametrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

const (
	counterMetricType = "counter"
	gaugeMetricType   = "gauge"
)

func TestServer_updateMetrics(t *testing.T) {
	storage, err := storage.NewStorage("mock")
	assert.NoError(t, err)
	s := &Server{
		router:  echo.New(),
		storage: storage,
		Log:     slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})),
	}
	s.storage.Open()
	defer s.storage.Close()
	s.ConfigureRouter()

	urlPath := "/update/:mtype/:mname/:mvalue"

	type want struct {
		code        int
		response    string
		contentType string
	}
	type params struct {
		mtype  string
		mname  string
		mvalue string
	}
	tests := []struct {
		name   string
		want   want
		target params
	}{
		{
			name: "add new counter",
			want: want{
				code:        http.StatusOK,
				response:    "updated",
				contentType: "text/plain; charset=UTF-8",
			},
			target: params{
				mtype:  counterMetricType,
				mname:  "c1",
				mvalue: "1",
			},
		},
		{
			name: "add new gauge",
			want: want{
				code:        http.StatusOK,
				response:    "updated",
				contentType: "text/plain; charset=UTF-8",
			},
			target: params{

				mtype:  gaugeMetricType,
				mname:  "g1",
				mvalue: "1.5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("Content-Type", "text/plain; charset=utf-8")
			rec := httptest.NewRecorder()
			c := s.router.NewContext(req, rec)
			c.SetPath(urlPath)
			c.SetParamNames("mtype", "mname", "mvalue")
			c.SetParamValues(tt.target.mtype, tt.target.mname, tt.target.mvalue)

			if err := s.updateMetrics(c); err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want.code, rec.Code)
			assert.Equal(t, tt.want.response, rec.Body.String())
			assert.Equal(t, tt.want.contentType, rec.Result().Header.Get("Content-Type"))
		})
	}
}

func TestServer_getMetricsByNames(t *testing.T) {
	storage, err := storage.NewStorage("mock")
	assert.NoError(t, err)
	s := &Server{
		router:  echo.New(),
		storage: storage,
		Log:     slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})),
	}
	s.storage.Open()
	defer s.storage.Close()
	s.ConfigureRouter()

	urlPath := "/value/:mtype/:mname"

	type want struct {
		code        int
		response    string
		contentType string
	}
	type params struct {
		mtype string
		mname string
	}
	tests := []struct {
		name   string
		want   want
		target params
	}{
		{
			name: "get counter metric, when is nill",
			want: want{
				code:        http.StatusNotFound,
				response:    "not found",
				contentType: "text/plain; charset=UTF-8",
			},
			target: params{
				mtype: "counter",
				mname: "c1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("Content-Type", "text/plain; charset=utf-8")
			rec := httptest.NewRecorder()
			c := s.router.NewContext(req, rec)
			c.SetPath(urlPath)
			c.SetParamNames("mtype", "mname")
			c.SetParamValues(tt.target.mtype, tt.target.mname)

			if err := s.getMetricsByNames(c); err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want.code, rec.Code)
			assert.Equal(t, tt.want.response, rec.Body.String())
			assert.Equal(t, tt.want.contentType, rec.Result().Header.Get("Content-Type"))
		})
	}
}
