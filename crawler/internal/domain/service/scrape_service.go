package service

type ScrapeService interface {
	ScrapeVideo(url string) (string, error)
}
