// Package logger provides logger.
package logger

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	so := StdOut{}

	tableTests := []struct {
		name       string
		env        string
		wantErr    error
		wantOut    []string
		notWantOut []string
	}{
		{
			name:    "Should Return Dev Logger",
			env:     "dev",
			wantErr: nil,
			wantOut: []string{
				`level=DEBUG`,
				`msg="debug layer"`,
				`msg="debug layer format Debug"`,
				`msg="debug layer with args"`,
				`argLevel=Debug`,
				`level=INFO`,
				`msg="info layer"`,
				`msg="info layer format Info"`,
				`msg="info layer with args"`,
				`argLevel=Info`,
				`level=WARN`,
				`msg="warn layer"`,
				`msg="warn layer format Warn"`,
				`msg="warn layer with args"`,
				`argLevel=Warn`,
				`level=ERROR`,
				`msg="error layer"`,
				`msg="error layer format Error"`,
				`msg="error layer with args"`,
				`argLevel=Error`,
			},
			notWantOut: []string{
				`service=test_service`,
			},
		},
		{
			name:    "Should Return Prod Logger",
			env:     "prod",
			wantErr: nil,
			wantOut: []string{
				`"service":"test_service"`,
				`"level":"INFO"`,
				`"msg":"info layer2"`,
				`"msg":"info layer format Info"`,
				`"msg":"info layer with args"`,
				`"argLevel":"Info"`,
				`"level":"WARN"`,
				`"msg":"warn layer"`,
				`"msg":"warn layer format Warn"`,
				`"msg":"warn layer with args"`,
				`"argLevel":"Warn"`,
				`"level":"ERROR"`,
				`"msg":"error layer"`,
				`"msg":"error layer format Error"`,
				`"msg":"error layer with args"`,
				`"argLevel":"Error"`,
			},
			notWantOut: []string{
				`"level":"DEBUG"`,
				`"msg":"debug layer"`,
				`"msg":"debug layer format Debug"`,
				`"msg":"debug layer with args"`,
				`"argLevel":"Debug"`,
			},
		},
		{
			name:    "Should Return Error",
			env:     "test",
			wantErr: ErrEnv,
		},
	}

	for _, tt := range tableTests {
		t.Run(tt.name, func(t *testing.T) {
			so.Acquire()
			logger, err := New(tt.env, "test_service")
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}

			assert.Implements(t, (*Logger)(nil), logger)
			logger.Debug("debug layer")
			logger.DebugF("debug layer format %v", "Debug")
			logger.DebugW("debug layer with args", "argLevel", "Debug")
			logger.Info("info layer")
			logger.InfoF("info layer format %v", "Info")
			logger.InfoW("info layer with args", "argLevel", "Info")
			logger.Warn("warn layer")
			logger.WarnF("warn layer format %v", "Warn")
			logger.WarnW("warn layer with args", "argLevel", "Warn")
			logger.Error("error layer")
			logger.ErrorF("error layer format %v", "Error")
			logger.ErrorW("error layer with args", "argLevel", "Error")
			out := so.Release()

			for _, wantOut := range tt.wantOut {
				assert.Contains(t, out, wantOut)
			}

			for _, notWantOut := range tt.notWantOut {
				assert.NotContains(t, out, notWantOut)
			}
		})
	}
}

func Test_Must(t *testing.T) {
	t.Run("Should Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Must(nil, errors.New(""))
		})
	})

	t.Run("Should Not Panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Must(&Log{}, nil)
		})
	})
}

// -------------
// --- Utils ---
// -------------

type StdOut struct {
	r      *os.File
	w      *os.File
	oldOut *os.File
}

func (s *StdOut) Acquire() {
	s.r, s.w, _ = os.Pipe()
	s.oldOut = os.Stdout
	os.Stdout = s.w
}

func (s *StdOut) Release() string {
	_ = s.w.Close()
	out, _ := io.ReadAll(s.r)
	os.Stdout = s.oldOut
	return string(out)
}
