package translation

import (
	"errors"
	"log"
	"nomenclature/internal/repository"
)

// TranslationService implements the fail-safe translation pipeline.
type TranslationService interface {
	Translate(req TranslationRequest) (TranslationResult, error)
}

type translationService struct {
	translationRepository repository.Querier
	ProviderPrimary       TranslationProvider
}

func NewTranslationService(provider TranslationProvider, translationRepository repository.Querier) TranslationService {
	return &translationService{
		ProviderPrimary:       provider,
		translationRepository: translationRepository,
	}
}

func (s *translationService) Translate(req TranslationRequest) (TranslationResult, error) {
	// Translation Memory Lookup
	// if res, found := s.translationRepository.Lookup(req); found {
	// 	return res, nil
	// }

	// Primary Provider
	res, err := s.ProviderPrimary.Translate(req)
	if err != nil {
		log.Printf("Primary provider failed: %v", err)
		return TranslationResult{Text: req.Text, Original: req.Text}, errors.New("translation failed")
	}

	// Store in Memory
	// s.translationRepository.Store(req, res)
	return res, nil
}
