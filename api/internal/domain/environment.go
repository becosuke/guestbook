package domain

import (
	"strings"

	"github.com/pkg/errors"
)

type Environment string

const (
	EnvUnknown     Environment = "unknown"
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
	EnvTest        Environment = "test"
)

var (
	ErrEnvUnknown = errors.New("unknown environment")
)

func NewEnvironment(s string) (Environment, error) {
	switch strings.ToLower(s) {
	case "development":
		return EnvDevelopment, nil
	case "production":
		return EnvProduction, nil
	case "test":
		return EnvTest, nil
	default:
		return EnvUnknown, ErrEnvUnknown
	}
}

func (e Environment) String() string {
	return string(e)
}
