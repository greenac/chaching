package env

type IBaseEnv interface {
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetString(key string) string
	IsSet(key string) bool
	AllSettings() map[string]interface{}
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

type Env struct {
	baseEnv IBaseEnv
}

func (e *Env) GetBool(key string) bool {
	return e.baseEnv.GetBool(key)
}

func (e *Env) GetFloat64(key string) float64 {
	return e.baseEnv.GetFloat64(key)
}

func (e *Env) GetInt(key string) int {
	return e.baseEnv.GetInt(key)
}

func (e *Env) GetString(key string) string {
	return e.baseEnv.GetString(key)
}

func (e *Env) IsSet(key string) bool {
	return e.baseEnv.IsSet(key)
}

func (e *Env) GetAll() map[string]interface{} {
	return e.baseEnv.AllSettings()
}
