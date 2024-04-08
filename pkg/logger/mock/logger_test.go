package loggermock

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/pass-keeper/pkg/logger"
)

func Test_New(t *testing.T) {
	so := StdOut{}

	tableTests := []struct {
		name       string
		env        string
		wantOut    []string
		notWantOut []string
	}{
		{
			name:       "Should Return Dev Logger",
			env:        "dev",
			wantOut:    []string{},
			notWantOut: []string{},
		},
		{
			name:       "Should Return Prod Logger",
			env:        "prod",
			wantOut:    []string{},
			notWantOut: []string{},
		},
		{
			name: "Should Return Error",
			env:  "test",
		},
	}

	for _, tt := range tableTests {
		t.Run(tt.name, func(t *testing.T) {
			so.Acquire()
			log, err := NewMock()
			if err != nil {
				return
			}

			assert.Implements(t, (*logger.Logger)(nil), log)
			log.Debug("debug layer")
			log.DebugF("debug layer format %v", "Debug")
			log.DebugW("debug layer with args", "argLevel", "Debug")
			log.Info("info layer")
			log.InfoF("info layer format %v", "Info")
			log.InfoW("info layer with args", "argLevel", "Info")
			log.Warn("warn layer")
			log.WarnF("warn layer format %v", "Warn")
			log.WarnW("warn layer with args", "argLevel", "Warn")
			log.Error("error layer")
			log.ErrorF("error layer format %v", "Error")
			log.ErrorW("error layer with args", "argLevel", "Error")
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
