package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/tenSunFree/travel_audio_guide_api/internal/delivery/http/response"
	"github.com/tenSunFree/travel_audio_guide_api/internal/domain"
	"github.com/tenSunFree/travel_audio_guide_api/internal/usecase"
)

// AudioHandler only knows the AudioUsecase interface; unaware of any infrastructure details
type AudioHandler struct {
	uc usecase.AudioUsecase
}

func NewAudioHandler(uc usecase.AudioUsecase) *AudioHandler {
	return &AudioHandler{uc: uc}
}

func (h *AudioHandler) GetAudio(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		response.Error(w, http.StatusNotFound, "Not Found")
		return
	}

	lang := parts[1]

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	result, err := h.uc.Execute(r.Context(), lang, page)
	if err != nil {
		// errors.As extracts AppError from the error chain — type-safe, no string parsing needed
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			response.Error(w, int(appErr.Code), appErr.Message)
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, result)
}
