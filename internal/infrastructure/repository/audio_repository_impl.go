package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tenSunFree/travel_audio_guide_api/internal/domain/entity"
	domainRepo "github.com/tenSunFree/travel_audio_guide_api/internal/domain/repository"
	"github.com/tenSunFree/travel_audio_guide_api/internal/infrastructure/client"
)

// audioRepositoryImpl implements domain.repository.AudioRepository
//
// Responsibilities:
//   - Call the client to obtain raw bytes
//   - Parse (map) the bytes into a domain entity
//   - Contains no business logic (language validation is not here)
type audioRepositoryImpl struct {
	client *client.TravelTaipeiClient
}

// NewAudioRepository returns a domain interface; callers are unaware of the concrete type
func NewAudioRepository(c *client.TravelTaipeiClient) domainRepo.AudioRepository {
	return &audioRepositoryImpl{client: c}
}

func (r *audioRepositoryImpl) GetAudio(ctx context.Context, lang string, page int) (*entity.AudioList, error) {
	body, err := r.client.FetchAudio(ctx, lang, page)
	if err != nil {
		return nil, err
	}

	// bytes → domain entity (mapper responsibility)
	var result entity.AudioList
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse upstream response: %w", err)
	}

	return &result, nil
}
