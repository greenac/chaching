package env

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type mockEnvBase struct {
	BoolVal   bool
	IntVal    int
	FloatVal  float64
	StringVal string
	IsSetVal  bool
	Settings  map[string]interface{}
}

func (me mockEnvBase) GetBool(key string) bool {
	return me.BoolVal
}

func (me mockEnvBase) GetFloat64(key string) float64 {
	return me.FloatVal
}

func (me mockEnvBase) GetInt(key string) int {
	return me.IntVal
}

func (me mockEnvBase) GetString(key string) string {
	return me.StringVal
}

func (me mockEnvBase) IsSet(key string) bool {
	return me.IsSetVal
}

func (me mockEnvBase) AllSettings() map[string]interface{} {
	return me.Settings
}

func TestEnv_GetBool(t *testing.T) {
	Convey("TestEnv_GetBool", t, func() {
		var val bool = true
		env := Env{baseEnv: mockEnvBase{BoolVal: val}}
		So(env.GetBool("something"), ShouldEqual, val)
	})
}

func TestEnv_GetFloat64(t *testing.T) {
	Convey("TestEnv_GetFloat64", t, func() {
		var val float64 = 1.1
		env := Env{baseEnv: mockEnvBase{FloatVal: val}}
		So(env.GetFloat64("something"), ShouldEqual, val)
	})
}

func TestEnv_GetInt(t *testing.T) {
	Convey("TestEnv_GetInt", t, func() {
		var val = 1
		env := Env{baseEnv: mockEnvBase{IntVal: val}}
		So(env.GetInt("something"), ShouldEqual, val)
	})
}

func TestEnv_GetString(t *testing.T) {
	Convey("TestEnv_GetString", t, func() {
		var val = "beach"
		env := Env{baseEnv: mockEnvBase{StringVal: val}}
		So(env.GetString("something"), ShouldEqual, val)
	})
}

func TestEnv_IsSet(t *testing.T) {
	Convey("TestEnv_GetString", t, func() {
		var val = false
		env := Env{baseEnv: mockEnvBase{IsSetVal: val}}
		So(env.IsSet("something"), ShouldEqual, val)
	})
}

func TestEnv_GetAll(t *testing.T) {
	Convey("TestEnv_GetAll", t, func() {
		var val = map[string]interface{}{"yippie": "yay"}
		env := Env{baseEnv: mockEnvBase{Settings: val}}
		So(env.GetAll(), ShouldResemble, val)
	})
}
