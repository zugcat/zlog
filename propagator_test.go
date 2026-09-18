package zlog_test

import (
	"bytes"
	"encoding/json/v2"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	"github.com/zugcat/zlog"
)

type headerReaderWriter map[string]*commonpb.Payload

var _ workflow.HeaderWriter = headerReaderWriter{}
var _ workflow.HeaderReader = headerReaderWriter{}

func (h headerReaderWriter) Set(key string, payload *commonpb.Payload) {
	h[key] = payload
}

func (h headerReaderWriter) Get(key string) (*commonpb.Payload, bool) {
	payload, ok := h[key]
	return payload, ok
}

func (h headerReaderWriter) ForEachKey(f func(string, *commonpb.Payload) error) error {
	for k, v := range h {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

func TestPropagator(t *testing.T) {
	ctx := t.Context()

	buf := &bytes.Buffer{}
	slog.SetDefault(zlog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{})))
	p := zlog.NewPropagator()

	t.Run("inject/extract", func(t *testing.T) {
		hdrs := make(headerReaderWriter)
		{
			// shadow the context here
			ctx := zlog.With(ctx, "key1", "value1", "key2", 2)
			ctx = zlog.With(zlog.WithGroup(ctx, "group"), "key1", "valueGroup")

			require.NoError(t, p.Inject(ctx, hdrs))
		}

		ctx, err := p.Extract(ctx, hdrs)
		require.NoError(t, err)

		slog.ErrorContext(ctx, "test")
		msg := make(map[string]any)

		require.NoError(t, json.Unmarshal(buf.Bytes(), &msg))
		t.Log(msg)

		assert.Equal(t, "value1", msg["key1"])
		assert.Equal(t, "valueGroup", msg["group"].(map[string]any)["key1"])
	})
	t.Run("workflow", func(t *testing.T) {
		hdrs := make(headerReaderWriter)

		{
			// shadow the context
			ctx := zlog.With(ctx, "key1", "workflowValue")
			require.NoError(t, p.Inject(ctx, hdrs))
		}

		s := &testsuite.WorkflowTestSuite{}
		env := s.NewTestWorkflowEnvironment()
		env.ExecuteWorkflow(func(ctx workflow.Context) error {
			ctx, err := p.ExtractToWorkflow(ctx, hdrs)
			require.NoError(t, err)

			hdrs2 := make(headerReaderWriter)
			err = p.InjectFromWorkflow(ctx, hdrs2)
			require.NoError(t, err)

			for k, v := range hdrs {
				assert.Equal(t, v, hdrs2[k])
			}

			return nil
		})
	})
}
