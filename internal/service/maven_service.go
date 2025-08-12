package service

import "cargo-m/internal/repository"

type MavenService struct {
	mavenRepo *repository.MavenRepo
}

func NewMavenService(mavenRepo *repository.MavenRepo) *MavenService {
	return &MavenService{mavenRepo: mavenRepo}
}
