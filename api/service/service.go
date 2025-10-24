package service

import (
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

type pingRes struct {
	Pong string `json:"pong"`
}

func (s *Service) Ping() (*pingRes, error) {
	return &pingRes{
		Pong: "pong",
	}, nil
}