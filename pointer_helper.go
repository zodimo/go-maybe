package maybe

import "reflect"

func FromPtrDereferenced[T any](ptr *T) Maybe[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

func FromPtr[T any](ptr *T) Maybe[*T] {
	if ptr == nil {
		return None[*T]()
	}
	return Some(ptr)
}

// for interface type, not any
func FromInterface[T any](interfaceValue T) Maybe[T] {

	value := reflect.ValueOf(&interfaceValue).Elem()
	kind := value.Kind()
	if kind == reflect.Interface {
		if value.IsNil() {
			return None[T]()
		}
	}
	return Some(interfaceValue)
}

func (m Maybe[T]) ToPtr() *T {
	if m.IsNone() {
		return nil
	}
	val := m.value
	return &val
}
