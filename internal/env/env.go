package env

import (
	"github.com/spf13/viper"
)

type GoEnv string

const (
	GoEnvLocal GoEnv = "local"
	GoEnvDev   GoEnv = "dev"
	GoEnvProd  GoEnv = "prod"
)

type IBaseEnv interface {
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetString(key string) string
	IsSet(key string) bool
	AllSettings() map[string]interface{}
	SetConfigName(in string)
	AddConfigPath(in string)
	SetConfigType(in string)
	SetConfigFile(in string)
	AutomaticEnv()
	ReadInConfig() error
}

type IEnv interface {
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetString(key string) string
	IsSet(key string) bool
	GetAll() map[string]interface{}
}

var _ IEnv = (*Env)(nil)

func NewEnv(filePath string, baseEnv IBaseEnv) (*Env, error) {
	ret := &Env{BaseEnv: baseEnv}
	if filePath == "" {
		baseEnv.SetConfigName(".env")
		baseEnv.AddConfigPath(".")
		baseEnv.SetConfigType("env")
	} else {
		baseEnv.SetConfigFile(filePath)
	}

	baseEnv.AutomaticEnv()
	if err := baseEnv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// EnvFile file not found; but ignored
		} else {
			return nil, err
		}
	}

	return ret, nil
}

type Env struct {
	BaseEnv IBaseEnv
}

func (e *Env) GetBool(key string) bool {
	return e.BaseEnv.GetBool(key)
}

func (e *Env) GetFloat64(key string) float64 {
	return e.BaseEnv.GetFloat64(key)
}

func (e *Env) GetInt(key string) int {
	return e.BaseEnv.GetInt(key)
}

func (e *Env) GetString(key string) string {
	return e.BaseEnv.GetString(key)
}

func (e *Env) IsSet(key string) bool {
	return e.BaseEnv.IsSet(key)
}

func (e *Env) GetAll() map[string]interface{} {
	return e.BaseEnv.AllSettings()
}
