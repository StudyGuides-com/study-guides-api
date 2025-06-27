package store

import (
	"github.com/lucsky/cuid"
)

func GetCUID() string {
	return cuid.New()
}
