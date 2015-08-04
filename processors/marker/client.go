package marker

import (
	"github.com/index0h/go-tracker/entities"
)

type Client struct {
	ClientID [16]byte

	Fields entities.Hash
}
