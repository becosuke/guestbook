package config

type Environment int

const (
	EnvUnknown Environment = iota
	EnvDevelopment
	EnvProduction
	EnvTest
)

func NewEnvironment(s string) Environment {
	e := EnvUnknown
	_ = e.Set(s)
	return e
}

func (e *Environment) Set(s string) error {
	switch s {
	case "Development", "DEVELOPMENT", "development":
		*e = EnvDevelopment
	case "Production", "PRODUCTION", "production":
		*e = EnvProduction
	case "Test", "TEST", "test":
		*e = EnvTest
	default:
		*e = EnvUnknown
	}
	return nil
}

func (e Environment) Get() interface{} {
	return e
}

func (e Environment) String() string {
	switch e {
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
