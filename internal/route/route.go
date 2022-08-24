package route

import (
	"go-wai-wong/internal/route/sumapi"

	"github.com/go-chi/chi"
)

func Install(r chi.Router) {
	sumapi.InstallRoutes(r)
}
