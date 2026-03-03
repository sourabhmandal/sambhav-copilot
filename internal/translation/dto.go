package translation

// TranslationRequest represents a request for translation.
type TranslationRequest struct {
	CompanyID      int64    `json:"company_id"`
	SourceLanguage string   `json:"source_language"`
	TargetLanguage string   `json:"target_language"`
	Text           []string `json:"text"`
}

// TranslationResult represents the result of a translation.
type TranslationResult struct {
	CompanyID      int64         `json:"company_id"`
	SourceLanguage string        `json:"source_language"`
	TargetLanguage string        `json:"target_language"`
	Text           []Translation `json:"text"`
}

type Translation struct {
	NormalizedHash string  `json:"hash"`
	Original       string  `json:"original"`
	Translated     string  `json:"translated"`
	Confidence     float64 `json:"confidence"`
	Provider       *string `json:"provider"`
}
