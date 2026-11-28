package shortenlinkservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/huyvo1112/golang_course/generator"
	"github.com/huyvo1112/golang_course/repository"
)

const defaultLength = 7

var ErrorDoesNotExist = errors.New("URL does not exist")

type ShortenLink interface {
	CreateShortLink(ctx context.Context, url string, expirationTime time.Duration) (string, error)
	GetURLFromCode(ctx context.Context, code string) (string, error)
}

type shortenLink struct {
	r       repository.URLStorage
	codeGen generator.CodeGenerator
}

func NewShortenLinkService(urlStorage repository.URLStorage, codeGen generator.CodeGenerator) ShortenLink {
	return &shortenLink{
		r:       urlStorage,
		codeGen: codeGen,
	}
}

func (s *shortenLink) CreateShortLink(ctx context.Context, url string, expirationTime time.Duration) (string, error) {
	// Generate code
	code, err := s.codeGen.GenerateCode(defaultLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate short code: %w", err)
	}
	// Store URL theo code
	err = s.r.SaveURL(ctx, code, url, expirationTime)
	if err != nil {
		return "", fmt.Errorf("failed to save shortened URL: %w", err)
	}
	return code, nil
}

func (s *shortenLink) GetURLFromCode(ctx context.Context, code string) (string, error) {
	res, err := s.r.GetURL(ctx, code)
	if errors.Is(err, repository.ErrorDoesNotExist) {
		return "", ErrorDoesNotExist
	}
	return res, err
}
