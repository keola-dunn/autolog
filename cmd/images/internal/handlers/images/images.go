package images

type ImagesHandler struct {
}

type ImagesHandlerConfig struct {
}

func NewAuthHandler(config ImagesHandlerConfig) (*ImagesHandler, error) {
	return &ImagesHandler{}, nil
}
