package common

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var ValidatorPool = &sync.Pool{
	New: func() interface{} { return validator.New() },
}
