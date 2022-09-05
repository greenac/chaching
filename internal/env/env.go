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
