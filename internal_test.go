package zlog

import (
	"encoding/json/v2"
	"log/slog"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type valuer int

func (v valuer) LogValue() slog.Value {
	return slog.IntValue(int(v))
}

func TestJSON(t *testing.T) {
	g := &handler{
		group: "group",
	}

	now := time.Now().UTC()

	a := &handler{
		attrs: []slog.Attr{
			{
				Key:   "stringKey",
				Value: slog.StringValue("stringValue"),
			},
			{
				Key:   "boolKey",
				Value: slog.BoolValue(true),
			},
			{
				Key:   "intKey",
				Value: slog.Int64Value(5),
			},
			{
				Key:   "uintKey",
				Value: slog.Uint64Value(6),
			},
			{
				Key:   "floatKey",
				Value: slog.Float64Value(4.5),
			},
			{
				Key:   "durationKey",
				Value: slog.DurationValue(5 * time.Second),
			},
			{
				Key:   "timeKey",
				Value: slog.TimeValue(now),
			},
			{
				Key:   "groupKey",
				Value: slog.GroupValue(slog.GroupAttrs("groupValue", slog.Attr{Key: "key", Value: slog.StringValue("value")})),
			},
			{
				Key:   "valuerKey",
				Value: slog.AnyValue(valuer(5)),
			},
		},
	}

	jg, err := json.Marshal(g)
	require.NoError(t, err)
	ja, err := json.Marshal(a)
	require.NoError(t, err)

	t.Log(string(jg))
	t.Log(string(ja))

	g2, a2 := &handler{}, &handler{}
	require.NoError(t, json.Unmarshal(jg, g2))
	require.NoError(t, json.Unmarshal(ja, a2))

	assert.Equal(t, g.group, g2.group)

	// instead of ElementsMatch, get out and walk because slog.Value
	// can decide to do delayed evaluation things
	assert.Equal(t, len(a.attrs), len(a2.attrs))
	for _, attr := range a.attrs {
		i := slices.IndexFunc(a2.attrs, func(attr2 slog.Attr) bool {
			return attr2.Key == attr.Key
		})
		require.GreaterOrEqual(t, i, 0)
		attr2 := a2.attrs[i]
		attr.Value = attr.Value.Resolve()
		attr2.Value = attr2.Value.Resolve()
		assert.Equal(t, attr.Value, attr2.Value)
	}
}
