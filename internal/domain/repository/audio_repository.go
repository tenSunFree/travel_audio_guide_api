package repository

import (
	"context"

	"github.com/tenSunFree/travel_audio_guide_api/internal/domain/entity"
)

// AudioRepository defines the data capabilities required by the domain.
//
// Dependency inversion key:
//   - Defined in the domain layer (inner layer)
//   - Implemented by the infrastructure layer (outer layer)
//   - Usecases depend only on this interface, unaware of any implementation details
type AudioRepository interface {
	GetAudio(ctx context.Context, lang string, page int) (*entity.AudioList, error)
}
