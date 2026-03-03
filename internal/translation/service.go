package translation

import (
	"context"
	"errors"
	"log"
	"math/big"
	"nomenclature/internal/repository"
	"nomenclature/internal/translationproviders"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type translationService struct {
	translationRepository repository.Querier
	translator            translationproviders.TranslationProvider
}

func NewTranslationService(provider translationproviders.TranslationProvider, translationRepository repository.Querier) TranslationService {
	return &translationService{
		translator:            provider,
		translationRepository: translationRepository,
	}
}

func (s *translationService) Translate(req TranslationRequest) (TranslationResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	all_hashes := make([]string, 0, len(req.Text))
	for _, text := range req.Text {
		hash := GenerateHash(req.CompanyID, req.SourceLanguage, req.TargetLanguage, text)
		all_hashes = append(all_hashes, hash)
	}
	// Translation Memory Lookup
	savedTranslationResponse, err := s.translationRepository.GetAllTranslationsByHashes(ctx, &repository.GetAllTranslationsByHashesParams{
		CompanyID: &req.CompanyID,
		Column2:   all_hashes,
	})
	var response = TranslationResult{
		CompanyID:      req.CompanyID,
		SourceLanguage: req.SourceLanguage,
		TargetLanguage: req.TargetLanguage,
		Text:           []Translation{},
	}
	if err == nil && len(savedTranslationResponse) > 0 {
		for _, translation := range savedTranslationResponse {
			log.Printf("CACHE HIT: Found translation in memory for hash %s: %s -> %s (Confidence: %f, Provider: %s)", translation.NormalizedHash, translation.OriginalText,
				translation.TranslatedText, float64(translation.ConfidenceScore.Int.Int64()), *translation.Provider)
			response.Text = append(response.Text, Translation{
				NormalizedHash: translation.NormalizedHash,
				Original:       translation.OriginalText,
				Translated:     translation.TranslatedText,
				Confidence:     float64(translation.ConfidenceScore.Int.Int64()),
				Provider:       nil,
			})
		}
	}

	nonCachedTranslations := make([]string, 0)
	for _, text := range req.Text {
		hash := GenerateHash(req.CompanyID, req.SourceLanguage, req.TargetLanguage, text)
		found := false
		for _, cachedTranslation := range savedTranslationResponse {
			if cachedTranslation.NormalizedHash == hash {
				found = true
				break
			}
		}
		if !found {
			nonCachedTranslations = append(nonCachedTranslations, text)
		}
	}

	if len(nonCachedTranslations) == 0 {
		return response, nil
	}

	log.Printf("CACHE MISS: No translation found in memory for %d texts. Calling translation provider...", len(nonCachedTranslations))

	// If not found in memory, call the translation provider
	translationResp, err := s.translator.Translate(translationproviders.TranslationRequest{
		CompanyID:      req.CompanyID,
		Text:           nonCachedTranslations,
		SourceLanguage: req.SourceLanguage,
		TargetLanguage: req.TargetLanguage,
	})
	if err != nil {
		log.Printf("Primary provider failed: %v", err)
		return response, errors.New("translation failed")
	}

	// Store in Memory
	var translationInputs []repository.TranslationInput
	for _, translation := range translationResp {
		translationInputs = append(translationInputs, repository.TranslationInput{
			CompanyID:       req.CompanyID,
			NormalizedHash:  GenerateHash(req.CompanyID, req.SourceLanguage, req.TargetLanguage, translation.Original),
			SourceLanguage:  req.SourceLanguage,
			TargetLanguage:  req.TargetLanguage,
			OriginalText:    translation.Original,
			TranslatedText:  translation.Translated,
			ConfidenceScore: translation.Confidence,
			Provider:        translation.Provider,
		})
	}

	var (
		n                            = len(translationResp)
		bulkInsertTranslationsParams repository.BulkInsertTranslationsParams
	)
	bulkInsertTranslationsParams = repository.BulkInsertTranslationsParams{
		Column1: make([]int64, n),
		Column2: make([]string, n),
		Column3: make([]string, n),
		Column4: make([]string, n),
		Column5: make([]string, n),
		Column6: make([]string, n),
		Column7: make([]pgtype.Numeric, n),
		Column8: make([]string, n),
	}
	for i, translation := range translationResp {
		bulkInsertTranslationsParams.Column1[i] = req.CompanyID
		bulkInsertTranslationsParams.Column2[i] = GenerateHash(req.CompanyID, req.SourceLanguage, req.TargetLanguage, translation.Original)
		bulkInsertTranslationsParams.Column3[i] = req.SourceLanguage
		bulkInsertTranslationsParams.Column4[i] = req.TargetLanguage
		bulkInsertTranslationsParams.Column5[i] = translation.Original
		bulkInsertTranslationsParams.Column6[i] = translation.Translated
		bulkInsertTranslationsParams.Column7[i] = pgtype.Numeric{Int: big.NewInt(int64(translation.Confidence)), Valid: true, Exp: -3}
		bulkInsertTranslationsParams.Column8[i] = translation.Provider
	}
	savedTranslationResponse, err = s.translationRepository.BulkInsertTranslations(ctx, &bulkInsertTranslationsParams)
	if err != nil {
		log.Printf("Failed to save translation: %v", err)
		// Even if saving fails, we can return the translation result
		// but without the confidence score and provider information.
		return response, nil
	}

	for _, translation := range savedTranslationResponse {
		response.Text = append(response.Text, Translation{
			NormalizedHash: translation.NormalizedHash,
			Original:       translation.OriginalText,
			Translated:     translation.TranslatedText,
			Confidence:     float64(translation.ConfidenceScore.Int.Int64()),
			Provider:       translation.Provider,
		})
	}

	return response, nil
}
