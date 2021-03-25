package shortner

import (
	"time"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid  = errors.New("Redirect Invalid")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

func NewRedirectService(redirectRepo RedirectRepository) RedirectService {
	return redirectService{
		redirectRepo: redirectRepo,
	}
}

func (r redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepo.Find(code)
}

func (r redirectService) Store(redirect *Redirect) error {
	err := validate.Validate(redirect)
	if err != nil {
		return errors.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}
