package app

import (
	"crypto/sha256"
	"fmt"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthorService struct {
	repo              repository.AuthorRepository
	siteConfigService *SiteConfigService
}

func NewAuthorService(repo repository.AuthorRepository, siteConfigService *SiteConfigService) *AuthorService {
	return &AuthorService{repo: repo, siteConfigService: siteConfigService}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *AuthorService) Create(name string, password string) (*model.Author, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(name, hash)
}

func (s *AuthorService) SetPassword(name string, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	author, err := s.repo.FindByName(name)
	if err != nil {
		return err
	}
	author.PasswordHash = hash
	err = s.repo.Update(author)
	return err

}

func (s *AuthorService) FindByName(name string) (*model.Author, error) {
	return s.repo.FindByName(name)
}

func (s *AuthorService) Authenticate(name string, password string) bool {
	author, err := s.repo.FindByName(name)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(author.PasswordHash), []byte(password))
	return err == nil
}

func (s *AuthorService) getSecretKey() string {
	siteConfig, err := s.siteConfigService.GetSiteConfig()
	if err != nil {
		panic(err)
	}
	if siteConfig.Secret == "" {
		siteConfig.Secret = RandStringRunes(64)
		err = s.siteConfigService.UpdateSiteConfig(siteConfig)
		if err != nil {
			panic(err)
		}
	}
	return siteConfig.Secret
}

func (s *AuthorService) CreateToken(name string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(name + s.getSecretKey()))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%x", name, hash.Sum(nil)), nil
}

func (s *AuthorService) ValidateToken(token string) (bool, string) {
	parts := strings.Split(token, ".")
	witness := parts[len(parts)-1]
	name := strings.Join(parts[:len(parts)-1], ".")

	hash := sha256.New()
	_, err := hash.Write([]byte(name + s.getSecretKey()))
	if err != nil {
		return false, ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil)) == witness, name
}
