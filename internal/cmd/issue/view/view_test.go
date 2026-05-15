package view

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdView_WebFlag(t *testing.T) {
	cmd := NewCmdView()

	flag := cmd.Flags().Lookup(flagWeb)
	if assert.NotNil(t, flag, "--web flag should be registered") {
		assert.Equal(t, "w", flag.Shorthand, "--web should have -w shorthand")
		assert.Equal(t, "false", flag.DefValue, "--web should default to false")
		assert.Equal(t, "bool", flag.Value.Type())
	}
}

func TestViewWeb(t *testing.T) {
	cases := []struct {
		name     string
		project  string
		server   string
		arg      string
		expected string
	}{
		{
			name:     "it builds URL from project and numeric key",
			project:  "TEST",
			server:   "https://example.atlassian.net",
			arg:      "42",
			expected: "https://example.atlassian.net/browse/TEST-42",
		},
		{
			name:     "it uses fully-qualified key as-is",
			project:  "TEST",
			server:   "https://example.atlassian.net",
			arg:      "OTHER-7",
			expected: "https://example.atlassian.net/browse/OTHER-7",
		},
		{
			name:     "it works with no configured project",
			project:  "",
			server:   "https://example.atlassian.net",
			arg:      "TEST-1",
			expected: "https://example.atlassian.net/browse/TEST-1",
		},
	}

	// `true` is a benign no-op on POSIX systems; this prevents browser.Browse
	// from launching a real browser during the test.
	t.Setenv("JIRA_BROWSER", "true")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(configProject, tc.project)
			viper.Set(configServer, tc.server)
			t.Cleanup(viper.Reset)

			out := captureStdout(t, func() {
				viewWeb([]string{tc.arg})
			})

			assert.Equal(t, tc.expected, strings.TrimSpace(out))
		})
	}
}

func TestViewWeb_HonorsBrowseServerOverride(t *testing.T) {
	t.Setenv("JIRA_BROWSER", "true")

	viper.Reset()
	viper.Set(configProject, "TEST")
	viper.Set(configServer, "https://api.example.com")
	viper.Set("browse_server", "https://web.example.com")
	t.Cleanup(viper.Reset)

	out := captureStdout(t, func() {
		viewWeb([]string{"1"})
	})

	assert.Equal(t, "https://web.example.com/browse/TEST-1", strings.TrimSpace(out))
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w

	done := make(chan struct{})
	var buf bytes.Buffer
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	fn()

	assert.NoError(t, w.Close())
	os.Stdout = orig
	<-done
	return buf.String()
}
