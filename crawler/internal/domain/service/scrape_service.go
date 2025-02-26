package service

import "context"

type ScrapeService interface {
	ScrapeVideo(ctx context.Context, url string) (string, error)
}
