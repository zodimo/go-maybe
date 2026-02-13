//go:build goexperiment.jsonv2

package maybe

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"reflect"
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
	var ptr *T
	if err := json.UnmarshalDecode(dec, &ptr); err != nil {
		return err
	}
	if ptr == nil {
		// Check if T is a pointer type
		var zero T
		if reflect.TypeOf(&zero).Elem().Kind() == reflect.Pointer {
			// For pointer types, null means Some(nil)
			m.hasValue = true
			m.value = zero
		} else {
			// For non-pointer types, null means None
			m.hasValue = false
			m.value = zero
		}
		return nil
	}
	m.hasValue = true
	m.value = *ptr
	return nil
}
