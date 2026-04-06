package usecase

import (
	"context"
	"fmt"

	"github.com/tenSunFree/travel_audio_guide_api/internal/domain"
	"github.com/tenSunFree/travel_audio_guide_api/internal/domain/entity"
	domainRepo "github.com/tenSunFree/travel_audio_guide_api/internal/domain/repository"
)

// validLangs holds supported language codes; this is a business rule owned by the usecase layer
var validLangs = map[string]bool{
	"zh-tw": true,
	"zh-cn": true,
	"en":    true,
	"ja":    true,
	"ko":    true,
}

// AudioUsecase interface
// The delivery layer depends on this interface, not on a concrete struct,
// making it easy to swap implementations or inject mocks for testing.
type AudioUsecase interface {
	Execute(ctx context.Context, lang string, page int) (*entity.AudioList, error)
}

type audioUsecase struct {
	repo domainRepo.AudioRepository
}

func NewAudioUsecase(repo domainRepo.AudioRepository) AudioUsecase {
	return &audioUsecase{repo: repo}
}

func (u *audioUsecase) Execute(ctx context.Context, lang string, page int) (*entity.AudioList, error) {
	// Business rule: validate language code
	if !validLangs[lang] {
		return nil, domain.NewBadRequest(fmt.Sprintf("unsupported language: %s", lang))
	}

	// Business rule: default page number
	if page < 1 {
		page = 1
	}

	result, err := u.repo.GetAudio(ctx, lang, page)
	if err != nil {
		return nil, domain.NewUpstreamFail("upstream API request failed: " + err.Error())
	}

	return result, nil
}
