package zlog_test

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zugcat/zlog"
)

func TestHandler(t *testing.T) {
	ctx := t.Context()
	buf := &bytes.Buffer{}

	slog.SetDefault(
		zlog.NewDefault(buf),
	)

	ctx = zlog.WithGroup(ctx, "mygroup")
	ctx = zlog.With(ctx, "attr1", "value1", slog.String("attr2", "value2"), 5, "x")

	slog.InfoContext(ctx, "log message")

	out := make(map[string]any)
	err := json.UnmarshalRead(buf, &out, jsontext.AllowDuplicateNames(true))
	require.NoError(t, err)

	assert.Equal(t, "value1", out["mygroup"].(map[string]any)["attr1"])
	assert.Equal(t, "value2", out["mygroup"].(map[string]any)["attr2"])
	// XXX: !badkey is not documented externally as the key
	assert.EqualValues(t, 5, out["mygroup"].(map[string]any)["!badkey"])
	// XXX: !badvalue is not documented externally as the key
	assert.Equal(t, "x", out["mygroup"].(map[string]any)["!badvalue"])
}
