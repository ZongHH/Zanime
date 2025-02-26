package service

import (
	"context"
)

type CrawlerService interface {
	Start(ctx context.Context)
}
