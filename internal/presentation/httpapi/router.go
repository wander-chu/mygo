package httpapi

import (
	"mygo/internal/option"

	"github.com/go-chi/chi"
)

func newRouter(_ *option.Options) chi.Router {
	router := chi.NewRouter()
	return router
}
