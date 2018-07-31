package reloggerZap

import (
	"github.com/remicro/trifle"
	"testing"
)

func (fx *fixture) Logger() *zapLogger {
	return &zapLogger{
		logger: fx.log,
	}
}

func TestZapLogger_Critical(t *testing.T) {
	fx := newFx(t)

	exp := trifle.String()
	result := map[string]interface{}{}

	fx.ws.EXPECT().WriteJson(t, &result)
	fx.ws.EXPECT().Sync().Return(nil)

	fx.Logger().
		Critical().
		Log(exp)
	assertLogEntry(t, "ERROR", exp, result)
}

func TestZapLogger_Debug(t *testing.T) {
	fx := newFx(t)

	exp := trifle.String()
	result := map[string]interface{}{}

	fx.ws.EXPECT().WriteJson(t, &result)
	fx.ws.EXPECT().Sync().Return(nil)

	fx.Logger().
		Debug().
		Log(exp)
	assertLogEntry(t, "DEBUG", exp, result)
}

func TestZapLogger_Error(t *testing.T) {
	fx := newFx(t)

	exp := trifle.String()
	result := map[string]interface{}{}

	fx.ws.EXPECT().WriteJson(t, &result)
	fx.ws.EXPECT().Sync().Return(nil)

	fx.Logger().
		Error().
		Log(exp)
	assertLogEntry(t, "ERROR", exp, result)
}

func TestZapLogger_Info(t *testing.T) {
	fx := newFx(t)

	exp := trifle.String()
	result := map[string]interface{}{}

	fx.ws.EXPECT().WriteJson(t, &result)
	fx.ws.EXPECT().Sync().Return(nil)

	fx.Logger().
		Info().
		Log(exp)
	assertLogEntry(t, "INFO", exp, result)
}

func TestZapLogger_Warn(t *testing.T) {
	fx := newFx(t)

	exp := trifle.String()
	result := map[string]interface{}{}

	fx.ws.EXPECT().WriteJson(t, &result)
	fx.ws.EXPECT().Sync().Return(nil)

	fx.Logger().
		Warn().
		Log(exp)
	assertLogEntry(t, "WARN", exp, result)
}

func TestNew(t *testing.T) {
	New()
}
