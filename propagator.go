package zlog

import (
	"context"
	"slices"

	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

const headerKey = "github.com/zugcat/zlog"

type propagator struct{}

var _ workflow.ContextPropagator = propagator{}

// NewPropagator creates a new [workflow.ContextPropagator] for zlog
func NewPropagator() workflow.ContextPropagator {
	return propagator{}
}

// Inject implements [workflow.ContextPropagator]
func (p propagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	return getHandler(ctx).inject(writer)
}

func (h *handler) inject(writer workflow.HeaderWriter) error {
	stack := []*handler{}

	for h != nil {
		stack = append(stack, h)
		h = h.up
	}

	slices.Reverse(stack)

	payload, err := converter.GetDefaultDataConverter().ToPayload(stack)
	if err != nil {
		return err
	}

	writer.Set(headerKey, payload)
	return nil
}

// Extract implements [workflow.ContextPropagator]
func (p propagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if payload, ok := reader.Get(headerKey); ok {
		var stack []*handler
		if err := converter.GetDefaultDataConverter().FromPayload(payload, &stack); err != nil {
			return ctx, err
		}

		h := getHandler(ctx)
		for _, s := range stack {
			if s.group != "" {
				h = h.WithGroup(s.group).(*handler)
			} else if len(s.attrs) > 0 {
				h = h.WithAttrs(s.attrs).(*handler)
			}
		}

		return withHandler(ctx, h), nil
	}
	return ctx, nil
}

// InjectFromWorkflow implements [workflow.ContextPropagator]
func (p propagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	stack, ok := ctx.Value(headerKey).([]*handler)
	if !ok {
		return nil
	}
	payload, err := converter.GetDefaultDataConverter().ToPayload(stack)
	if err != nil {
		return err
	}

	writer.Set(headerKey, payload)
	return nil
}

// ExtractToWorkflow implements [workflow.ContextPropagator]
func (p propagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if payload, ok := reader.Get(headerKey); ok {
		var stack []*handler
		if err := converter.GetDefaultDataConverter().FromPayload(payload, &stack); err != nil {
			return ctx, err
		}
		return workflow.WithValue(ctx, headerKey, stack), nil
	}

	return ctx, nil
}
