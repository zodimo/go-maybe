package maybe

func Map[T any, R any](m Maybe[T], f func(T) R) Maybe[R] {
	if m.IsSome() {
		return Some(f(m.value))
	}
	return None[R]()
}
func FlatMap[T any, R any](m Maybe[T], f func(T) Maybe[R]) Maybe[R] {
	if m.IsSome() {
		return f(m.value)
	}
	return None[R]()
}
