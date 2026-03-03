package repository

type TranslationInput struct {
	CompanyID       int64
	NormalizedHash  string
	SourceLanguage  string
	TargetLanguage  string
	OriginalText    string
	TranslatedText  string
	ConfidenceScore float64
	Provider        string
}
