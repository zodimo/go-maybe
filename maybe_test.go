package maybe

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewMaybe(t *testing.T) {
	t.Run("creates empty maybe for int", func(t *testing.T) {
		m := NewMaybe[int]()
		if m.IsSome() {
			t.Error("NewMaybe should create None value")
		}
		if !m.IsNone() {
			t.Error("NewMaybe should create None value")
		}
	})

	t.Run("creates empty maybe for string", func(t *testing.T) {
		m := NewMaybe[string]()
		if m.IsSome() {
			t.Error("NewMaybe should create None value")
		}
		if !m.IsNone() {
			t.Error("NewMaybe should create None value")
		}
	})
}

func TestSome(t *testing.T) {
	t.Run("creates Some with int value", func(t *testing.T) {
		m := Some(42)
		if !m.IsSome() {
			t.Error("Some should create Some value")
		}
		if m.IsNone() {
			t.Error("Some should not create None value")
		}

		value, err := m.Unwrap()
		if err != nil {
			t.Errorf("Some should unwrap without error, got: %v", err)
		}
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
	})

	t.Run("creates Some with string value", func(t *testing.T) {
		m := Some("hello")
		if !m.IsSome() {
			t.Error("Some should create Some value")
		}

		value, err := m.Unwrap()
		if err != nil {
			t.Errorf("Some should unwrap without error, got: %v", err)
		}
		if value != "hello" {
			t.Errorf("Expected 'hello', got %v", value)
		}
	})

	t.Run("creates Some with zero value", func(t *testing.T) {
		m := Some(0)
		if !m.IsSome() {
			t.Error("Some should create Some value even with zero value")
		}

		value, err := m.Unwrap()
		if err != nil {
			t.Errorf("Some should unwrap without error, got: %v", err)
		}
		if value != 0 {
			t.Errorf("Expected 0, got %v", value)
		}
	})
}

func TestNone(t *testing.T) {
	t.Run("creates None for int", func(t *testing.T) {
		m := None[int]()
		if m.IsSome() {
			t.Error("None should create None value")
		}
		if !m.IsNone() {
			t.Error("None should create None value")
		}
	})

	t.Run("creates None for string", func(t *testing.T) {
		m := None[string]()
		if m.IsSome() {
			t.Error("None should create None value")
		}
		if !m.IsNone() {
			t.Error("None should create None value")
		}
	})
}

func TestIsSome(t *testing.T) {
	t.Run("returns true for Some", func(t *testing.T) {
		m := Some(42)
		if !m.IsSome() {
			t.Error("IsSome should return true for Some value")
		}
	})

	t.Run("returns false for None", func(t *testing.T) {
		m := None[int]()
		if m.IsSome() {
			t.Error("IsSome should return false for None value")
		}
	})
}

func TestIsNone(t *testing.T) {
	t.Run("returns false for Some", func(t *testing.T) {
		m := Some(42)
		if m.IsNone() {
			t.Error("IsNone should return false for Some value")
		}
	})

	t.Run("returns true for None", func(t *testing.T) {
		m := None[int]()
		if !m.IsNone() {
			t.Error("IsNone should return true for None value")
		}
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("returns value and nil error for Some", func(t *testing.T) {
		m := Some(42)
		value, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error for Some, got: %v", err)
		}
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
	})

	t.Run("returns zero value and error for None", func(t *testing.T) {
		m := None[int]()
		value, err := m.Unwrap()
		if err == nil {
			t.Error("Unwrap should return error for None")
		}
		if err.Error() != "none" {
			t.Errorf("Expected error 'none', got: %v", err.Error())
		}
		if value != 0 {
			t.Errorf("Expected zero value 0, got %v", value)
		}
	})

	t.Run("returns zero value and error for None string", func(t *testing.T) {
		m := None[string]()
		value, err := m.Unwrap()
		if err == nil {
			t.Error("Unwrap should return error for None")
		}
		if value != "" {
			t.Errorf("Expected empty string, got %v", value)
		}
	})
}

func TestUnwrapUnsafe(t *testing.T) {
	t.Run("returns value for Some", func(t *testing.T) {
		m := Some(42)
		value := m.UnwrapUnsafe()
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
	})

	t.Run("panics for None", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("UnwrapUnsafe should panic for None")
			}
		}()

		m := None[int]()
		m.UnwrapUnsafe()
	})
}

func TestUnwrapOr(t *testing.T) {
	t.Run("returns value for Some", func(t *testing.T) {
		m := Some(42)
		value := m.UnwrapOr(100)
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
	})

	t.Run("returns default value for None", func(t *testing.T) {
		m := None[int]()
		value := m.UnwrapOr(100)
		if value != 100 {
			t.Errorf("Expected 100, got %v", value)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		m := Some("hello")
		value := m.UnwrapOr("default")
		if value != "hello" {
			t.Errorf("Expected 'hello', got %v", value)
		}

		m2 := None[string]()
		value2 := m2.UnwrapOr("default")
		if value2 != "default" {
			t.Errorf("Expected 'default', got %v", value2)
		}
	})

	t.Run("works with zero values", func(t *testing.T) {
		m := Some(0)
		value := m.UnwrapOr(100)
		if value != 0 {
			t.Errorf("Expected 0, got %v", value)
		}

		m2 := None[int]()
		value2 := m2.UnwrapOr(0)
		if value2 != 0 {
			t.Errorf("Expected 0, got %v", value2)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("applies function to Some value", func(t *testing.T) {
		m := Some(5)
		result := m.Map(func(x int) int { return x * 2 })

		if !result.IsSome() {
			t.Error("Map should return Some when applied to Some")
		}

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != 10 {
			t.Errorf("Expected 10, got %v", value)
		}
	})

	t.Run("returns None when applied to None", func(t *testing.T) {
		m := None[int]()
		result := m.Map(func(x int) int { return x * 2 })

		if !result.IsNone() {
			t.Error("Map should return None when applied to None")
		}
	})

	t.Run("applies string transformation", func(t *testing.T) {
		m := Some("hello")
		result := m.Map(func(s string) string { return s + " world" })

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != "hello world" {
			t.Errorf("Expected 'hello world', got %v", value)
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("applies function to Some value", func(t *testing.T) {
		m := Some(5)
		result := m.FlatMap(func(x int) Maybe[int] {
			if x > 0 {
				return Some(x * 2)
			}
			return None[int]()
		})

		if !result.IsSome() {
			t.Error("FlatMap should return Some when function returns Some")
		}

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("FlatMap result should unwrap without error, got: %v", err)
		}
		if value != 10 {
			t.Errorf("Expected 10, got %v", value)
		}
	})

	t.Run("returns None when function returns None", func(t *testing.T) {
		m := Some(-5)
		result := m.FlatMap(func(x int) Maybe[int] {
			if x > 0 {
				return Some(x * 2)
			}
			return None[int]()
		})

		if !result.IsNone() {
			t.Error("FlatMap should return None when function returns None")
		}
	})

	t.Run("returns None when applied to None", func(t *testing.T) {
		m := None[int]()
		result := m.FlatMap(func(x int) Maybe[int] {
			return Some(x * 2)
		})

		if !result.IsNone() {
			t.Error("FlatMap should return None when applied to None")
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("returns Some when predicate is true", func(t *testing.T) {
		m := Some(10)
		result := m.Filter(func(x int) bool { return x > 5 })

		if !result.IsSome() {
			t.Error("Filter should return Some when predicate is true")
		}

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Filter result should unwrap without error, got: %v", err)
		}
		if value != 10 {
			t.Errorf("Expected 10, got %v", value)
		}
	})

	t.Run("returns None when predicate is false", func(t *testing.T) {
		m := Some(3)
		result := m.Filter(func(x int) bool { return x > 5 })

		if !result.IsNone() {
			t.Error("Filter should return None when predicate is false")
		}
	})

	t.Run("returns None when applied to None", func(t *testing.T) {
		m := None[int]()
		result := m.Filter(func(x int) bool { return true })

		if !result.IsNone() {
			t.Error("Filter should return None when applied to None")
		}
	})
}

func TestOrElse(t *testing.T) {
	t.Run("returns original value for Some", func(t *testing.T) {
		m := Some(42)
		result := m.OrElse(100)

		if result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("returns else value for None", func(t *testing.T) {
		m := None[int]()
		result := m.OrElse(100)

		if result != 100 {
			t.Errorf("Expected 100, got %v", result)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		m := None[string]()
		result := m.OrElse("default")

		if result != "default" {
			t.Errorf("Expected 'default', got %v", result)
		}
	})
}

func TestOrElseGet(t *testing.T) {
	t.Run("returns original value for Some", func(t *testing.T) {
		m := Some(42)
		result := m.OrElseGet(func() int { return 100 })

		if result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("calls function and returns result for None", func(t *testing.T) {
		m := None[int]()
		called := false
		result := m.OrElseGet(func() int {
			called = true
			return 100
		})

		if !called {
			t.Error("Function should be called for None")
		}
		if result != 100 {
			t.Errorf("Expected 100, got %v", result)
		}
	})

	t.Run("does not call function for Some", func(t *testing.T) {
		m := Some(42)
		called := false
		result := m.OrElseGet(func() int {
			called = true
			return 100
		})

		if called {
			t.Error("Function should not be called for Some")
		}
		if result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	})
}

func TestOrElseError(t *testing.T) {
	t.Run("returns value and nil error for Some", func(t *testing.T) {
		m := Some(42)
		customErr := errors.New("custom error")
		value, err := m.OrElseError(customErr)

		if err != nil {
			t.Errorf("OrElseError should not return error for Some, got: %v", err)
		}
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
	})

	t.Run("returns zero value and custom error for None", func(t *testing.T) {
		m := None[int]()
		customErr := errors.New("custom error")
		value, err := m.OrElseError(customErr)

		if err != customErr {
			t.Errorf("Expected custom error, got: %v", err)
		}
		if value != 0 {
			t.Errorf("Expected zero value 0, got %v", value)
		}
	})

	t.Run("works with strings", func(t *testing.T) {
		m := None[string]()
		customErr := errors.New("string error")
		value, err := m.OrElseError(customErr)

		if err != customErr {
			t.Errorf("Expected custom error, got: %v", err)
		}
		if value != "" {
			t.Errorf("Expected empty string, got %v", value)
		}
	})
}

// Test chaining operations
func TestChaining(t *testing.T) {
	t.Run("chain map operations", func(t *testing.T) {
		m := Some(5)
		result1 := m.Map(func(x int) int { return x * 2 })
		result2 := result1.Map(func(x int) int { return x + 1 })

		value, err := result2.Unwrap()
		if err != nil {
			t.Errorf("Chained operations should unwrap without error, got: %v", err)
		}
		if value != 11 {
			t.Errorf("Expected 11, got %v", value)
		}
	})

	t.Run("chain with filter", func(t *testing.T) {
		m := Some(10)
		result1 := m.Map(func(x int) int { return x * 2 })
		result2 := result1.Filter(func(x int) bool { return x > 15 })

		value := result2.OrElse(0)
		if value != 20 {
			t.Errorf("Expected 20, got %v", value)
		}
	})

	t.Run("chain breaks on None", func(t *testing.T) {
		m := None[int]()
		result1 := m.Map(func(x int) int { return x * 2 })
		result2 := result1.Map(func(x int) int { return x + 1 })

		if !result2.IsNone() {
			t.Error("Chain should remain None when starting with None")
		}
	})
}

func TestNewMaybeEqualsNone(t *testing.T) {
	t.Run("new(Maybe[string]) == None[string]()", func(t *testing.T) {
		m := *new(Maybe[string])
		none := None[string]()

		if m != none {
			t.Error("new(Maybe[string]) should equal None[string]()")
		}

		// Also verify they behave the same
		if m.IsSome() != none.IsSome() {
			t.Error("new(Maybe[string]) and None[string]() should have same IsSome() result")
		}
		if m.IsNone() != none.IsNone() {
			t.Error("new(Maybe[string]) and None[string]() should have same IsNone() result")
		}
	})
}

func TestHelperMap(t *testing.T) {
	t.Run("applies function to Some value with type transformation", func(t *testing.T) {
		m := Some(5)
		result := Map(m, func(x int) string {
			return fmt.Sprintf("value: %d", x)
		})

		if !result.IsSome() {
			t.Error("Map should return Some when applied to Some")
		}

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != "value: 5" {
			t.Errorf("Expected 'value: 5', got %v", value)
		}
	})

	t.Run("returns None when applied to None", func(t *testing.T) {
		m := None[int]()
		result := Map(m, func(x int) string {
			return "transformed"
		})

		if !result.IsNone() {
			t.Error("Map should return None when applied to None")
		}
	})

	t.Run("transforms int to string", func(t *testing.T) {
		m := Some(42)
		result := Map(m, func(x int) string {
			return "number: 42"
		})

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != "number: 42" {
			t.Errorf("Expected 'number: 42', got %v", value)
		}
	})

	t.Run("transforms string to int", func(t *testing.T) {
		m := Some("123")
		result := Map(m, func(s string) int {
			return len(s)
		})

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != 3 {
			t.Errorf("Expected 3, got %v", value)
		}
	})

	t.Run("transforms int to float", func(t *testing.T) {
		m := Some(5)
		result := Map(m, func(x int) float64 {
			return float64(x) * 1.5
		})

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != 7.5 {
			t.Errorf("Expected 7.5, got %v", value)
		}
	})
}

func TestHelperFlatMap(t *testing.T) {
	t.Run("applies function to Some value with type transformation", func(t *testing.T) {
		m := Some(5)
		result := FlatMap(m, func(x int) Maybe[string] {
			if x > 0 {
				return Some("positive")
			}
			return None[string]()
		})

		if !result.IsSome() {
			t.Error("FlatMap should return Some when function returns Some")
		}

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("FlatMap result should unwrap without error, got: %v", err)
		}
		if value != "positive" {
			t.Errorf("Expected 'positive', got %v", value)
		}
	})

	t.Run("returns None when function returns None", func(t *testing.T) {
		m := Some(-5)
		result := FlatMap(m, func(x int) Maybe[string] {
			if x > 0 {
				return Some("positive")
			}
			return None[string]()
		})

		if !result.IsNone() {
			t.Error("FlatMap should return None when function returns None")
		}
	})

	t.Run("returns None when applied to None", func(t *testing.T) {
		m := None[int]()
		result := FlatMap(m, func(x int) Maybe[string] {
			return Some("transformed")
		})

		if !result.IsNone() {
			t.Error("FlatMap should return None when applied to None")
		}
	})

	t.Run("transforms int to string Maybe", func(t *testing.T) {
		m := Some(42)
		result := FlatMap(m, func(x int) Maybe[string] {
			return Some("number: 42")
		})

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("FlatMap result should unwrap without error, got: %v", err)
		}
		if value != "number: 42" {
			t.Errorf("Expected 'number: 42', got %v", value)
		}
	})

	t.Run("transforms string to int Maybe", func(t *testing.T) {
		m := Some("hello")
		result := FlatMap(m, func(s string) Maybe[int] {
			if len(s) > 0 {
				return Some(len(s))
			}
			return None[int]()
		})

		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("FlatMap result should unwrap without error, got: %v", err)
		}
		if value != 5 {
			t.Errorf("Expected 5, got %v", value)
		}
	})

	t.Run("conditional transformation returns None", func(t *testing.T) {
		m := Some("")
		result := FlatMap(m, func(s string) Maybe[int] {
			if len(s) > 0 {
				return Some(len(s))
			}
			return None[int]()
		})

		if !result.IsNone() {
			t.Error("FlatMap should return None when condition fails")
		}
	})
}

func TestNoneWithInterface(t *testing.T) {
	t.Run("creates None for interface{} without panic", func(t *testing.T) {
		m := None[interface{}]()
		if m.IsSome() {
			t.Error("None should create None value")
		}
		if !m.IsNone() {
			t.Error("None should create None value")
		}

		// Test unwrap
		value, err := m.Unwrap()
		if err == nil {
			t.Error("Unwrap should return error for None")
		}
		if value != nil {
			t.Errorf("Expected nil for interface{} zero value, got %v", value)
		}
	})

	t.Run("creates Some for interface{}", func(t *testing.T) {
		m := Some[interface{}]("test")
		if !m.IsSome() {
			t.Error("Some should create Some value")
		}
		value, err := m.Unwrap()
		if err != nil {
			t.Errorf("Some should unwrap without error, got: %v", err)
		}
		if value != "test" {
			t.Errorf("Expected 'test', got %v", value)
		}
	})
}
