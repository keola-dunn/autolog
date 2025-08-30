package random

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type ServiceIface interface {
	RandomUUID() (string, error)
	RandomString(int64) string
	RandomUpperAlphanumericString(length int64) string
}

func NewService() *Service {
	return &Service{}
}

type Service struct {
}

func (s *Service) RandomUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate random uuid: %w", err)
	}

	return id.String(), nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (s *Service) RandomString(length int64) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var upperAlphanumericRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func (s *Service) RandomUpperAlphanumericString(length int64) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(upperAlphanumericRunes))]
	}
	return string(b)
}
