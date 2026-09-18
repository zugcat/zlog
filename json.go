package zlog

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"log/slog"
	"time"
)

var _ json.MarshalerTo = handler{}
var _ json.UnmarshalerFrom = &handler{}

func unmarshalAttrFrom(d *jsontext.Decoder, attr *slog.Attr) error {
	proxy := &struct {
		Key   string     `json:"key"`
		Value slog.Value `json:"value"`
	}{}

	if err := json.UnmarshalDecode(d, proxy,
		json.WithUnmarshalers(
			json.UnmarshalFromFunc(unmarshalValueFrom),
		),
	); err != nil {
		return err
	}

	attr.Key = proxy.Key
	attr.Value = proxy.Value

	return nil
}

func unmarshalDurationFrom(d *jsontext.Decoder, value *time.Duration) (err error) {
	var str string
	if err := json.UnmarshalDecode(d, &str); err != nil {
		return err
	}

	*value, err = time.ParseDuration(str)
	return err
}

func unmarshalTimeFrom(d *jsontext.Decoder, value *time.Time) (err error) {
	var str string
	if err := json.UnmarshalDecode(d, &str); err != nil {
		return err
	}

	*value, err = time.Parse(time.RFC3339Nano, str)
	return err
}

var (
	boolKind     = slog.KindBool.String()
	durationKind = slog.KindDuration.String()
	timeKind     = slog.KindTime.String()
	float64Kind  = slog.KindFloat64.String()
	int64Kind    = slog.KindInt64.String()
	uint64Kind   = slog.KindUint64.String()
	groupKind    = slog.KindGroup.String()
)

func unmarshalValueFrom(d *jsontext.Decoder, value *slog.Value) error {
	proxy := &struct {
		Kind  string         `json:"kind"`
		Value jsontext.Value `json:"value"`
	}{}

	if err := json.UnmarshalDecode(d, proxy); err != nil {
		return err
	}

	var v any

	switch proxy.Kind {
	case boolKind:
		v = false

	case durationKind:
		v = time.Duration(0)

	case timeKind:
		v = time.Time{}

	case float64Kind:
		v = float64(0)

	case int64Kind:
		v = int64(0)

	case uint64Kind:
		v = uint64(0)

	case groupKind:
		v = []slog.Attr{}

	default:
		v = ""
	}

	if err := json.Unmarshal(proxy.Value, &v,
		json.WithUnmarshalers(
			json.JoinUnmarshalers(
				json.UnmarshalFromFunc(unmarshalAttrFrom),
				json.UnmarshalFromFunc(unmarshalDurationFrom),
				json.UnmarshalFromFunc(unmarshalTimeFrom),
			),
		),
	); err != nil {
		return err
	}
	*value = slog.AnyValue(v)

	return nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom], used for context propagation
func (h *handler) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	proxy := &struct {
		Attrs []slog.Attr `json:"attrs"`
		Group string      `json:"group"`
	}{}

	if err := json.UnmarshalDecode(d, proxy,
		json.WithUnmarshalers(
			json.UnmarshalFromFunc(unmarshalAttrFrom),
		),
	); err != nil {
		return err
	}

	h.attrs = proxy.Attrs
	h.group = proxy.Group

	return nil
}

func marshalAttrTo(e *jsontext.Encoder, attr slog.Attr) error {
	proxy := &struct {
		Key   string     `json:"key"`
		Value slog.Value `json:"value"`
	}{
		Key:   attr.Key,
		Value: attr.Value,
	}

	return json.MarshalEncode(e, proxy,
		json.WithMarshalers(
			json.MarshalToFunc(marshalValueTo),
		),
	)
}

func marshalDurationTo(e *jsontext.Encoder, value time.Duration) error {
	return json.MarshalEncode(e, value.String())
}

func marshalTimeTo(e *jsontext.Encoder, value time.Time) error {
	return json.MarshalEncode(e, value.Format(time.RFC3339Nano))
}

func marshalValueTo(e *jsontext.Encoder, value slog.Value) error {
	if value.Kind() == slog.KindLogValuer {
		value = value.Resolve()
	}

	proxy := &struct {
		Kind  string `json:"kind"`
		Value any    `json:"value"`
	}{
		Kind:  value.Kind().String(),
		Value: value.Any(),
	}

	return json.MarshalEncode(e, proxy,
		json.WithMarshalers(
			json.JoinMarshalers(
				json.MarshalToFunc(marshalAttrTo),
				json.MarshalToFunc(marshalDurationTo),
				json.MarshalToFunc(marshalTimeTo),
			),
		),
	)
}

// MarshalJSONTo implements [json.MarhsalerTo], used for context propagation
func (h handler) MarshalJSONTo(e *jsontext.Encoder) error {
	proxy := &struct {
		Attrs []slog.Attr `json:"attrs,omitempty"`
		Group string      `json:"group,omitempty"`
	}{
		Attrs: h.attrs,
		Group: h.group,
	}
	return json.MarshalEncode(e, proxy,
		json.WithMarshalers(
			json.MarshalToFunc(marshalAttrTo),
		),
	)
}
