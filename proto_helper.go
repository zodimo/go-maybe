package maybe

// From Proto to Maybe
func FromPtr[T any](ptr *T) Maybe[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

// From Maybe to Proto
func (m Maybe[T]) ToPtr() *T {
	if m.IsNone() {
		return nil
	}
	val := m.value
	return &val
}
