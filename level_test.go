package zlog_test

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"expvar"
	"flag"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zugcat/zlog"
)

func TestLevel(t *testing.T) {
	buf := &bytes.Buffer{}

	// don't let jsontext.Decoder get fancy on us
	r := struct{ io.Reader }{buf}

	dec := jsontext.NewDecoder(r)
	logger := zlog.NewDefault(buf)

	t.Logf("current log level is %s", zlog.Level())
	zlog.Level().Set("INFO")
	assert.Equal(t, "INFO", expvar.Get("logLevel").String())

	logger.Debug("debug level")
	logger.Info("info level")

	// only the info record should have made it
	data := make(map[string]any)
	require.NoError(t, json.UnmarshalDecode(dec, &data))
	assert.Equal(t, "info level", data["msg"])

	zlog.Level().Set("DEBUG")
	logger.Debug("debug level")

	clear(data)
	require.NoError(t, json.UnmarshalDecode(dec, &data))
	assert.Equal(t, "debug level", data["msg"])
}

func TestFlag(t *testing.T) {
	fs := flag.NewFlagSet("testing", flag.ContinueOnError)
	fs.Var(zlog.Level(), "loglevel", "sets logging level")

	err := fs.Parse([]string{"-loglevel=WARN"})
	require.NoError(t, err)

	assert.Equal(t, slog.LevelWarn, zlog.Level().Level())
	assert.Equal(t, slog.LevelWarn, zlog.Level().Get())
}
