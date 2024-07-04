package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"what-to.com/internal/config"
)

type EntityMiddleware struct {
	appRouter *mux.Router
	appConfig *config.Config
}

func NewEntityMiddleware(appConfig *config.Config, router *mux.Router) *EntityMiddleware {
	m := &EntityMiddleware{}
	m.appRouter = router
	m.appConfig = appConfig
	return m
}

func (m *EntityMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.appConfig.GetLogger().Info(fmt.Sprintf("Start  [HTTP] [%s] handler called from: %s, path: %s", r.Method, r.RemoteAddr, r.URL.Path))
		next.ServeHTTP(w, r)
		m.appConfig.GetLogger().Info(fmt.Sprintf("Finish [HTTP] [%s] handler called from: %s, path: %s", r.Method, r.RemoteAddr, r.URL.Path))
	})
}

func (m *EntityMiddleware) GetHandler() http.Handler {
	return m.Logging(m.appRouter)
}

func (m *EntityMiddleware) GetMiddlewareRouter() *mux.Router {
	return m.appRouter
}
