package repositories

import "glorified_hashmap/src/domain/models"

type RedirectRepository interface {
	Find(code string) (*models.Redirect, error)
	Store(redirect *models.Redirect) error
}
