package translationproviders

// TranslationProvider defines the interface for translation providers.
type TranslationProvider interface {
	Translate(req TranslationRequest) ([]Translation, error)
}
