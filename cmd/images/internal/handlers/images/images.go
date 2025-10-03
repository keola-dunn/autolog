package images

import "github.com/keola-dunn/autolog/internal/service/image"

type ImagesHandler struct {
	imageSvc image.ServiceIface
}

type ImagesHandlerConfig struct {
	ImageService image.ServiceIface
}

func NewHandler(config ImagesHandlerConfig) (*ImagesHandler, error) {
	return &ImagesHandler{
		imageSvc: config.ImageService,
	}, nil
}
