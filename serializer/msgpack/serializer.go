package msgpack

import (
	"github.com/hex_url_shortner/shortner"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

type Redirect struct{}

func (r Redirect) Decode(input []byte) (*shortner.Redirect, error) {
	redirect := shortner.Redirect{}
	err := msgpack.Unmarshal(input, &redirect)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return &redirect, nil
}

func (r Redirect) Encode(input *shortner.Redirect) ([]byte, error) {
	data, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return data, nil
}
