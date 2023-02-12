package config

import (
	"github.com/pkg/errors"
	"strings"
)

type Environment int

const (
	EnvUnknown Environment = iota
	EnvDevelopment
	EnvProduction
	EnvTest
)

var (
	ErrEnvUnknown = errors.New("unknown environment")
)

func NewEnvironment(s string) (Environment, error) {
	e := EnvUnknown
	err := e.Set(s)
	return e, err
}

func (e *Environment) Set(s string) error {
	switch strings.ToLower(s) {
	case "development":
		*e = EnvDevelopment
	case "production":
		*e = EnvProduction
	case "test":
		*e = EnvTest
	default:
		return ErrEnvUnknown
	}
	return nil
}

func (e *Environment) String() string {
	switch *e {
	case EnvDevelopment:
		return "development"
	case EnvProduction:
		return "production"
	case EnvTest:
		return "test"
	default:
		return "unknown"
	}
}
