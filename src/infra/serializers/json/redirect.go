package serializers_json

import (
	"encoding/json"
	"glorified_hashmap/src/domain/models"

	"github.com/pkg/errors"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*models.Redirect, error) {
	redirect := &models.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializers_json.redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *models.Redirect) ([]byte, error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializers_json.redirect.Encode")
	}
	return bytes, nil
}
