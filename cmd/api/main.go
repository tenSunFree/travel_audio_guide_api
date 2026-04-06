package main

import (
	"log"
	"net/http"

	"github.com/tenSunFree/travel_audio_guide_api/config"
	httpDelivery "github.com/tenSunFree/travel_audio_guide_api/internal/delivery/http"
	"github.com/tenSunFree/travel_audio_guide_api/internal/delivery/http/handler"
	infraClient "github.com/tenSunFree/travel_audio_guide_api/internal/infrastructure/client"
	infraRepo "github.com/tenSunFree/travel_audio_guide_api/internal/infrastructure/repository"
	"github.com/tenSunFree/travel_audio_guide_api/internal/usecase"
)

// main — the only place that knows about all layers.
//
// Dependency injection order (inner → outer):
//  1. config
//  2. infrastructure (client → repository impl)
//  3. usecase (injected with repository interface)
//  4. delivery (injected with usecase interface)
//  5. start server
func main() {
	// 1. Config
	cfg := config.Load()

	// 2. Infrastructure layer
	travelClient := infraClient.NewTravelTaipeiClient(cfg.UpstreamBaseURL, cfg.HTTPTimeout)
	audioRepo := infraRepo.NewAudioRepository(travelClient) // returns domain interface

	// 3. Usecase layer (only knows the domain interface, unaware of infraRepo)
	audioUsecase := usecase.NewAudioUsecase(audioRepo) // returns usecase interface

	// 4. Delivery layer (only knows the usecase interface)
	audioHandler := handler.NewAudioHandler(audioUsecase)
	swaggerHandler := handler.NewSwaggerHandler()
	router := httpDelivery.NewRouter(audioHandler, swaggerHandler)

	// 5. Start server
	addr := ":" + cfg.Port
	log.Printf("Travel Audio Guide API starting...")
	log.Printf("API:          http://localhost%s/open-api/{lang}/Media/Audio", addr)
	log.Printf("Swagger UI:   http://localhost%s/open-api/swagger/ui/index", addr)
	log.Printf("Swagger JSON: http://localhost%s/open-api/swagger/docs/V1", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
