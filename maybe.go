package maybe

import (
	"errors"
)

type Maybe[T any] struct {
	value    T
	hasValue bool
}

func NewMaybe[T any]() Maybe[T] {
	return None[T]()
}

func Some[T any](value T) Maybe[T] {
	return Maybe[T]{
		value:    value,
		hasValue: true,
	}
}
func None[T any]() Maybe[T] {
	var zero T
	return Maybe[T]{
		value:    zero,
		hasValue: false,
	}
}

func (m Maybe[T]) IsSome() bool {
	return m.hasValue
}
func (m Maybe[T]) IsNone() bool {
	return !m.hasValue
}
func (m Maybe[T]) Unwrap() (T, error) {
	if m.hasValue {
		return m.value, nil
	}
	var zero T
	return zero, errors.New("none")
}
func (m Maybe[T]) UnwrapUnsafe() T {
	value, err := m.Unwrap()
	if err != nil {
		panic(err)
	}
	return value
}

func (m Maybe[T]) UnwrapOr(defaultValue T) T {
	return m.OrElse(defaultValue)
}

func (m Maybe[T]) Map(f func(T) T) Maybe[T] {
	if m.hasValue {
		return Some(f(m.value))
	}
	return None[T]()
}

func (m Maybe[T]) FlatMap(f func(T) Maybe[T]) Maybe[T] {
	if m.hasValue {
		return f(m.value)
	}
	return None[T]()
}

func (m Maybe[T]) Filter(f func(T) bool) Maybe[T] {
	if m.hasValue && f(m.value) {
		return Some(m.value)
	}
	return None[T]()
}
func (m Maybe[T]) OrElse(elseValue T) T {
	if m.hasValue {
		return m.value
	}
	return elseValue
}
func (m Maybe[T]) OrElseGet(f func() T) T {
	if m.hasValue {
		return m.value
	}
	return f()
}
func (m Maybe[T]) OrElseError(err error) (T, error) {
	if m.hasValue {
		return m.value, nil
	}
	var zero T
	return zero, err
}
