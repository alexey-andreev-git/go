package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"what-to.com/internal/config"
	"what-to.com/internal/cors"
)

type EntityMiddleware struct {
	appRouter    *mux.Router
	appConfig    *config.Config
	handlerFuncs []func(http.Handler) http.Handler
}

func NewEntityMiddleware(appConfig *config.Config, router *mux.Router) *EntityMiddleware {
	m := &EntityMiddleware{}
	m.appRouter = router
	m.appConfig = appConfig
	m.handlerFuncs = []func(http.Handler) http.Handler{
		m.Cors,
		m.Logging,
		m.Recover,
	}
	return m
}

func (m *EntityMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msgfmt := "%s [HTTP] [%s] handler called from: %s, path: %s"
		m.appConfig.GetLogger().Info(fmt.Sprintf(msgfmt, "Start", r.Method, r.RemoteAddr, r.URL.Path))
		defer func() {
			m.appConfig.GetLogger().Info(fmt.Sprintf(msgfmt, "Finish", r.Method, r.RemoteAddr, r.URL.Path))
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *EntityMiddleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.appConfig.GetLogger().Error("Panic occurred: ", err.(error))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *EntityMiddleware) Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors := cors.NewDefaultCors()
		for k, v := range cors.GetHeaders() {
			w.Header().Set(k, v)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *EntityMiddleware) Use() {
	for _, f := range m.handlerFuncs {
		m.appRouter.Use(f)
	}
}

func (m *EntityMiddleware) ChainMiddleware() http.Handler {
	handler := http.Handler(m.GetMiddlewareRouter())
	for _, f := range m.handlerFuncs {
		handler = f(handler)
	}
	return handler
}

func (m *EntityMiddleware) GetMiddlewareRouter() *mux.Router {
	return m.appRouter
}
