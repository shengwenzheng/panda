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

type swzRes struct {
	Swz string `json:"swz"`
}

func (s *Service) Ping(_ *gin.Context) (*pingRes, error) {
	return &pingRes{
		Pong: "pong",
	}, nil
}

func (s *Service) Swz(_ *gin.Context) (*swzRes, error) {
	return &swzRes{
		Swz: "这是一段数据.........",
	}, nil
}
