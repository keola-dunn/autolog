package image

import (
	"context"
	"crypto/sha256"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"time"
)

type Image struct {
	Image image.Image

	id     string
	UserId string
	Title  string
	Path   string
	width  int64
	height int64
	SizeKb int64
	hash   string

	createdAt time.Time
	updatedAt time.Time
}

func (i *Image) Id() string {
	return i.id
}

func (i *Image) Width() int64 {
	return i.width
}

func (i *Image) Height() int64 {
	return i.height
}

func (i *Image) CreatedAt() time.Time {
	return i.createdAt
}

func (i *Image) UpdatedAt() time.Time {
	return i.updatedAt
}

func (s *Service) SaveImage(ctx context.Context, i Image) (*Image, error) {
	imageId, err := s.randomGenerator.RandomUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to create image id: %w", err)
	}

	i.id = imageId

	filePath := fmt.Sprintf("%s/%s.jpg", s.imagePrefix, imageId)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	i.Path = filePath

	defer file.Close()

	hash := sha256.New()

	if err = jpeg.Encode(io.MultiWriter(file, hash), i.Image, nil); err != nil {
		return nil, fmt.Errorf("failed to write image to file: %w", err)
	}

	fileStats, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	i.SizeKb = fileStats.Size() / 1000

	hashValue := hash.Sum(nil)
	i.hash = string(hashValue)

	i.width = int64(i.Image.Bounds().Dx())
	i.height = int64(i.Image.Bounds().Dy())

	query := `
	INSERT INTO images(id, user_id, title, path, width, height, imageSizeKb, hash)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)`

	if _, err := s.db.Exec(ctx, query, i.id, i.UserId, i.Title, i.Path,
		i.width, i.height, i.SizeKb, i.hash); err != nil {
		return nil, fmt.Errorf("failed to exec insert image query: %w", err)
	}

	return &i, nil
}
