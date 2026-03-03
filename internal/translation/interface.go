package translation

import "github.com/gin-gonic/gin"

type TranslationHandler interface {
	TranslateHandler(c *gin.Context)
}

type TranslationService interface {
	Translate(req TranslationRequest) (TranslationResult, error)
}
