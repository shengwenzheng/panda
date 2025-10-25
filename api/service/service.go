package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

type pingRes struct {
	Pong string `json:"pong"`
}

func (s *Service) Ping(_ *gin.Context) (*pingRes, error) {
	return &pingRes{
		Pong: "pong",
	}, nil
}