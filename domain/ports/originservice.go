package ports

import (
	"github.com/cgxarrie-go/prq/utils"
)

type OriginSvc interface {
	GetPRsURL(origin utils.Origin, status PRStatus) (url string, 
		err error)
	CreatePRsURL(origin utils.Origin) (url string, err error)
	PRLink(origin utils.Origin, id, text string) (url string, err error)
}