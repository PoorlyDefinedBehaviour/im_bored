package serializers

import "glorified_hashmap/src/domain/models"

type RedirectSerializer interface {
	Decode(input []byte) (*models.Redirect, error)
	Encode(input *models.Redirect) ([]byte, error)
}
