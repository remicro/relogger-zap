package reloggerZap

import (
	"math/rand"
	"testing"
	"time"

	"fmt"
	"github.com/remicro/relogger-zap/mock"
	"github.com/remicro/trifle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type fixture struct {
	log *zap.Logger
	ws  *mockWriteSyncer.MockWriteSyncer
}

func newFx(_ *testing.T) *fixture {
	ws := mockWriteSyncer.New()
	log := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			ws,
			zapcore.DebugLevel,
		),
	)
	return &fixture{
		ws:  ws,
		log: log,
	}
}

func TestEntry_Bool(t *testing.T) {
	t.Run("expect add boolean true field", func(t *testing.T) {
		e := newEntry(nil, zapcore.DebugLevel, nil)
		key := trifle.String()
		e.Bool(key, true)
		assert.Contains(t, e.fields, key)
		assert.True(t, e.fields[key].Equals(zap.Bool(key, true)))
	})
	t.Run("expect add boolean false field", func(t *testing.T) {
		e := newEntry(nil, zapcore.DebugLevel, nil)
		key := trifle.String()
		e.Bool(key, false)
		assert.Contains(t, e.fields, key)
		assert.True(t, e.fields[key].Equals(zap.Bool(key, false)))
	})
}

func TestEntry_String(t *testing.T) {
	t.Run("add empty valued string field", func(t *testing.T) {
		key, value := trifle.String(), trifle.StringN(0)
		e := newEntry(nil, zapcore.DebugLevel, nil)
		e.String(key, value)
		assert.Contains(t, e.fields, key)
		assert.True(t, e.fields[key].Equals(zap.String(key, value)))
	})
	t.Run("add valued string field", func(t *testing.T) {
		key, value := trifle.String(), trifle.String()
		e := newEntry(nil, zapcore.DebugLevel, nil)
		e.String(key, value)
		assert.Contains(t, e.fields, key)
		assert.True(t, e.fields[key].Equals(zap.String(key, value)))
	})
}

func TestEntry_Duration(t *testing.T) {
	t.Run("expect add field with empty duration", func(t *testing.T) {
		key := trifle.String()
		value := time.Duration(0)
		e := newEntry(nil, zapcore.DebugLevel, nil)
		e.Duration(key, value)
		assert.Contains(t, e.fields, key)
		assert.True(t, e.fields[key].Equals(zap.Duration(key, value)))
	})
	t.Run("expect add field with non-empty duration", func(t *testing.T) {
		key := trifle.String()
		value, err := time.ParseDuration("1s")
		require.NoError(t, err)
		e := newEntry(nil, zapcore.DebugLevel, nil)
		e.Duration(key, value)
		assert.Contains(t, e.fields, key)
		assert.True(t, e.fields[key].Equals(zap.Duration(key, value)))
	})
}

func TestEntry_Err(t *testing.T) {
	t.Run("expect add field with error", func(t *testing.T) {
		err := trifle.UnexpectedError()
		e := newEntry(nil, zapcore.DebugLevel, nil)
		e.Err(err)
		require.Contains(t, e.fields, "error")
		assert.True(t, e.fields["error"].Equals(zap.Error(err)))
	})
	t.Run("expect add field without error", func(t *testing.T) {
		e := newEntry(nil, zapcore.DebugLevel, nil)
		e.Err(nil)
		require.Contains(t, e.fields, "error")
		assert.True(t, e.fields["error"].Equals(zap.Error(nil)))
	})
}

func TestEntry_Time(t *testing.T) {
	e := newEntry(nil, zapcore.DebugLevel, nil)
	exp := time.Now()
	key := trifle.String()

	e.Time(key, exp)

	require.Contains(t, e.fields, key)
	assert.True(t, e.fields[key].Equals(zap.Time(key, exp)))
}

func TestEntry_Float64(t *testing.T) {
	e := newEntry(nil, zapcore.DebugLevel, nil)
	exp := rand.Float64()
	key := trifle.String()

	e.Float64(key, exp)

	require.Contains(t, e.fields, key)
	assert.True(t, e.fields[key].Equals(zap.Float64(key, exp)))
}

func TestEntry_Int(t *testing.T) {
	e := newEntry(nil, zapcore.DebugLevel, nil)
	exp := rand.Int()
	key := trifle.String()

	e.Int(key, exp)

	require.Contains(t, e.fields, key)
	assert.True(t, e.fields[key].Equals(zap.Int(key, exp)))
}

func TestEntry_Uint64(t *testing.T) {
	e := newEntry(nil, zapcore.DebugLevel, nil)
	exp := rand.Uint64()
	key := trifle.String()

	e.Uint64(key, exp)

	require.Contains(t, e.fields, key)
	assert.True(t, e.fields[key].Equals(zap.Uint64(key, exp)))
}
func assertLogEntry(t *testing.T, level, msg string, rec map[string]interface{}) {
	assert.Contains(t, rec, "T")
	assert.Contains(t, rec, "L")
	assert.Contains(t, rec, "M")
	assert.Equal(t, rec["L"], level)
	assert.Equal(t, rec["M"], msg)
}

func TestEntry_Log(t *testing.T) {

	t.Run("expect synced debug message", func(t *testing.T) {
		fx := newFx(t)
		e := newEntry(fx.log, zapcore.DebugLevel, nil)

		msg := trifle.String()
		saved := map[string]interface{}{}

		fx.ws.EXPECT().WriteJson(t, &saved)
		fx.ws.EXPECT().Sync().Return(nil)

		e.Log(msg)

		assertLogEntry(t, "DEBUG", msg, saved)
	})
	t.Run("expect synced error message", func(t *testing.T) {
		fx := newFx(t)
		e := newEntry(fx.log, zapcore.ErrorLevel, nil)

		msg := trifle.String()
		saved := map[string]interface{}{}

		fx.ws.EXPECT().WriteJson(t, &saved)
		fx.ws.EXPECT().Sync().Return(nil)

		e.Log(msg)

		assertLogEntry(t, "ERROR", msg, saved)
	})

	t.Run("expect synced info message", func(t *testing.T) {
		fx := newFx(t)
		e := newEntry(fx.log, zapcore.InfoLevel, nil)

		msg := trifle.String()
		saved := map[string]interface{}{}

		fx.ws.EXPECT().WriteJson(t, &saved)
		fx.ws.EXPECT().Sync().Return(nil)

		e.Log(msg)

		assertLogEntry(t, "INFO", msg, saved)
	})

	t.Run("expect synced warning message", func(t *testing.T) {
		fx := newFx(t)
		e := newEntry(fx.log, zapcore.WarnLevel, nil)

		msg := trifle.String()
		saved := map[string]interface{}{}

		fx.ws.EXPECT().WriteJson(t, &saved)
		fx.ws.EXPECT().Sync().Return(nil)

		e.Log(msg)

		assertLogEntry(t, "WARN", msg, saved)
	})
	t.Run("with field", func(t *testing.T) {
		t.Run("debug level", func(t *testing.T) {
			fx := newFx(t)
			e := newEntry(fx.log, zapcore.DebugLevel, nil)
			result := map[string]interface{}{}
			msg := trifle.String()

			key, value := trifle.StringN(5), trifle.String()
			fx.ws.EXPECT().WriteJson(t, &result)
			fx.ws.EXPECT().Sync().Return(nil)

			e.String(key, value).
				Log(msg)

			assertLogEntry(t, "DEBUG", msg, result)
			assert.Contains(t, result, key)
			assert.Equal(t, result[key], value)
		})

		t.Run("error level", func(t *testing.T) {
			fx := newFx(t)
			e := newEntry(fx.log, zapcore.ErrorLevel, nil)
			result := map[string]interface{}{}
			msg := trifle.String()

			key, value := trifle.StringN(5), trifle.String()
			fx.ws.EXPECT().WriteJson(t, &result)
			fx.ws.EXPECT().Sync().Return(nil)

			e.String(key, value).
				Log(msg)

			assertLogEntry(t, "ERROR", msg, result)
			assert.Contains(t, result, key)
			assert.Equal(t, result[key], value)
		})

		t.Run("info level", func(t *testing.T) {
			fx := newFx(t)
			e := newEntry(fx.log, zapcore.InfoLevel, nil)
			result := map[string]interface{}{}
			msg := trifle.String()

			key, value := trifle.StringN(5), trifle.String()
			fx.ws.EXPECT().WriteJson(t, &result)
			fx.ws.EXPECT().Sync().Return(nil)

			e.String(key, value).
				Log(msg)

			assertLogEntry(t, "INFO", msg, result)
			assert.Contains(t, result, key)
			assert.Equal(t, result[key], value)
		})

		t.Run("warn level", func(t *testing.T) {
			fx := newFx(t)
			e := newEntry(fx.log, zapcore.WarnLevel, nil)
			result := map[string]interface{}{}
			msg := trifle.String()

			key, value := trifle.StringN(5), trifle.String()
			fx.ws.EXPECT().WriteJson(t, &result)
			fx.ws.EXPECT().Sync().Return(nil)

			e.String(key, value).
				Log(msg)

			assertLogEntry(t, "WARN", msg, result)
			assert.Contains(t, result, key)
			assert.Equal(t, result[key], value)
		})
	})
}

func TestEntry_Logf(t *testing.T) {
	fx := newFx(t)
	e := newEntry(fx.log, zapcore.WarnLevel, nil)

	result := map[string]interface{}{}

	value := trifle.StringN(15)
	pattern := "test: %s"

	exp := fmt.Sprintf(pattern, value)
	fx.ws.EXPECT().WriteJson(t, &result)
	fx.ws.EXPECT().Sync().Return(nil)

	e.Logf(pattern, value)

	assertLogEntry(t, "WARN", exp, result)
}
