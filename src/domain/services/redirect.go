package services

import (
	"glorified_hashmap/src/domain/models"
	"glorified_hashmap/src/domain/repositories"
	"time"

	errors "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	RedirectNotFoundError = errors.New("Redirect Not Found")
	InvalidRedirectError  = errors.New("Invalid Redirect")
)

type RedirectService interface {
	Find(code string) (*models.Redirect, error)
	Store(redirect *models.Redirect) error
}

type redirectService struct {
	redirectRepository repositories.RedirectRepository
}

func New(redirectRepository repositories.RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepository,
	}
}

func (r *redirectService) Find(code string) (*models.Redirect, error) {
	return r.redirectRepository.Find(code)
}

func (r *redirectService) Store(redirect *models.Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errors.Wrap(InvalidRedirectError, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepository.Store(redirect)
}
