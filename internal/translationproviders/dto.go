package translationproviders

// TranslationRequest represents a request for translation.
type TranslationRequest struct {
	CompanyID      int64    `json:"company_id"`
	SourceLanguage string   `json:"source_language"`
	TargetLanguage string   `json:"target_language"`
	Text           []string `json:"text"`
}

type Translation struct {
	NormalizedHash string  `json:"hash"`
	Original       string  `json:"original"`
	Translated     string  `json:"translated"`
	Confidence     float64 `json:"confidence"`
	Provider       string  `json:"provider"`
}

type GoogleTranslateResponse struct {
	Translations []struct {
		TranslatedText string `json:"translatedText"`
	} `json:"translations"`
}
