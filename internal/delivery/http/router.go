package http

import (
	"net/http"
	"strings"

	"github.com/tenSunFree/travel_audio_guide_api/internal/delivery/http/handler"
	"github.com/tenSunFree/travel_audio_guide_api/internal/delivery/http/middleware"
)

// NewRouter wires routes and middleware.
// Accepts only handlers; unaware of any usecase or infrastructure details.
func NewRouter(audio *handler.AudioHandler, swagger *handler.SwaggerHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.Split(path, "/")

		// GET /open-api/{lang}/Media/Audio
		if len(parts) == 4 &&
			parts[0] == "open-api" &&
			strings.EqualFold(parts[2], "Media") &&
			strings.EqualFold(parts[3], "Audio") {
			audio.GetAudio(w, r)
			return
		}

		// Swagger UI
		if path == "" || path == "open-api/swagger/ui" || path == "open-api/swagger/ui/index" {
			swagger.UI(w, r)
			return
		}

		// Swagger Spec
		if path == "open-api/swagger/docs" || path == "open-api/swagger/docs/V1" {
			swagger.Spec(w, r)
			return
		}

		http.NotFound(w, r)
	})

	// middleware applied from innermost to outermost
	return middleware.Recover(
		middleware.Logging(
			middleware.CORS(mux),
		),
	)
}
