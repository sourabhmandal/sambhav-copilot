package translationproviders

import (
	"context"
	"fmt"
	"time"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"google.golang.org/api/option"
)

// TODO: implement OAuth2 for Google API calls in future iterations. Google Translate requires OAuth2 token
// GoogleTranslateProvider implements TranslationProvider for Google Cloud Translate.
type googleTranslateProvider struct {
	ProjectID string
}

func NewGoogleTranslateProvider(projectID string) TranslationProvider {
	return &googleTranslateProvider{
		ProjectID: projectID, // set this in ENV
	}
}

func (g *googleTranslateProvider) Translate(req TranslationRequest) ([]Translation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client, err := translate.NewTranslationClient(ctx, option.WithCredentialsFile("service-account.json"))
	if err != nil {
		return []Translation{}, err
	}
	defer client.Close()

	treq := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", g.ProjectID),
		Contents:           req.Text,
		TargetLanguageCode: req.TargetLanguage,
	}

	resp, err := client.TranslateText(ctx, treq)
	if err != nil {
		return []Translation{}, err
	}
	if len(resp.Translations) == 0 {
		return []Translation{}, fmt.Errorf("no translation returned")
	}
	var translations []Translation
	for i, tr := range resp.Translations {
		translations = append(translations, Translation{
			Provider:   "google",
			Original:   req.Text[i],
			Translated: tr.TranslatedText,
			Confidence: 90,
		})
	}
	return translations, nil
}
