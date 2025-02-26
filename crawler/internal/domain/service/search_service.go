package service

import "context"

type SearchService interface {
	SearchAnime(ctx context.Context, name, release, area, episode string) (string, error)
}
