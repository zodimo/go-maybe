//go:build goexperiment.jsonv2

package maybe

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// MarshalJSONTo implements jsontext.MarshalerTo interface for maximum performance.
// Writes directly to the encoder without intermediate buffer allocations.
func (m Maybe[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !m.hasValue {
		// When used with omitzero, this shouldn't be called for None values.
		// If called directly, output nothing.
		return nil
	}
	return json.MarshalEncode(enc, m.value)
}

// UnmarshalJSONFrom implements jsontext.UnmarshalerFrom interface.
// Sets hasValue=true when any JSON token is encountered for this field.
func (m *Maybe[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	m.hasValue = true
	return json.UnmarshalDecode(dec, &m.value)
}
