package model

type ShortenLink struct {
	ID        int    `json:"id"`
	ShortLink string `json:"short_link"`
	LongLink  string `json:"long_link"`
}

type ShortenLinkRequest struct {
	LongLink string `json:"long_link" binding:"required"`
}

type ShortenLinkResponse struct {
	ShortLink string `json:"short_link"`
}
