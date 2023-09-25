package errors

import (
	"fmt"

	"github.com/cgxarrie-go/prq/internal/utils"
)

type ErrInvalidRepositoryType struct {
	BaseError
}

func NewErrInvalidRepositoryType(origin utils.Remote) error {
	return ErrInvalidRepositoryType{
		BaseError: BaseError{
			message: fmt.Sprintf("invalid repository type %s", origin),
		},
	}
}