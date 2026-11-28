package shortenlinkhandler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	shortenlinkservice "github.com/huyvo1112/golang_course/service/shorten_link_service"
)

type ShortenLink interface {
	CreateShortenLink(c *gin.Context)
	RedirectShortenLink(c *gin.Context)
}

type shortenLink struct {
	svc shortenlinkservice.ShortenLink
}

func NewShortenLink(svc shortenlinkservice.ShortenLink) ShortenLink {
	return &shortenLink{svc: svc}
}

type shortenLinkInput struct {
	Exp int64  `json:"exp"`
	URL string `json:"url" binding:"required,url"`
}

// @Summary Create Shorten Link
// @Description Create Shorten Link
// @Tags Links
// @Produce application/json
// @Success 200 {object} string "Shorten URL generated successfully!"
// @Failure 500 {object} string "Internal Server Error"
// @Router /v1/links/shorten [post]
func (s *shortenLink) CreateShortenLink(c *gin.Context) {
	// Get input
	input := &shortenLinkInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call service
	code, err := s.svc.CreateShortLink(c, input.URL, time.Duration(input.Exp)*time.Second)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": code, "message": "Shorten URL generated successfully!"})
}

// Redirect Forward the request to the original url
// @Tags Links
// @Accept application/json
// @Produce application/json
// @Param code path string true "Shorten code"
// @Success 302
// @Router/v1/links/redirect/{code} [get]
func (s *shortenLink) RedirectShortenLink(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	url, err := s.svc.GetURLFromCode(c, code)
	if err != nil {
		if errors.Is(err, shortenlinkservice.ErrorDoesNotExist) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Code not exist"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, url)

}
