package services

import "os"

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

func (p *PostService) GetStaticPath(isUrl bool) string {
	staticPath := os.Getenv("STATIC_PATH") + "/posts/"
	if !isUrl {
		staticPath = "./" + staticPath
	}
	return staticPath
}

func (p *PostService) PostFileUpload() {

}
