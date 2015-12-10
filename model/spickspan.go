package model

import (
	"strconv"
)

type Service struct {
	Id   string `json:"id"`
	IP   string `json:"ip"`
	Port int `json:"port"`
}

func (s *Service) GetHttpUrl() string {
	return "http://" + s.IP + ":" + strconv.Itoa(s.Port)
}

type Provider interface {
	GetService(key string) (Service, error)
	Detect() bool
	Release(Service) error
	Init()
}

