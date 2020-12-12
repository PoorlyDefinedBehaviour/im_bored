package serializers_msgpack

import (
	"glorified_hashmap/src/domain/models"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*models.Redirect, error) {
	redirect := &models.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializers_msgpack.redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *models.Redirect) ([]byte, error) {
	bytes, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializers_msgpack.redirect.Encode")
	}
	return bytes, nil
}
