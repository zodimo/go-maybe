package maybe

import (
	"fmt"
	"testing"
)

// Tests for the generic type-transforming methods Map[R], FlatMap[R], and
// Match[R], plus the helper functions' short-circuit (non-invocation) behavior.

func TestMethodMapTypeChange(t *testing.T) {
	t.Run("transforms int to string via method", func(t *testing.T) {
		m := Some(42)
		result := m.Map(func(x int) string { return fmt.Sprintf("number: %d", x) })

		if !result.IsSome() {
			t.Error("Map should return Some when transforming to string")
		}
		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != "number: 42" {
			t.Errorf("Expected 'number: 42', got %v", value)
		}
	})

	t.Run("transforms string to int via method", func(t *testing.T) {
		m := Some("hello")
		result := m.Map(func(s string) int { return len(s) })

		if !result.IsSome() {
			t.Error("Map should return Some when transforming to int")
		}
		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != 5 {
			t.Errorf("Expected 5, got %v", value)
		}
	})

	t.Run("transforms int to float64 via method", func(t *testing.T) {
		m := Some(5)
		result := m.Map(func(x int) float64 { return float64(x) * 1.5 })

		if !result.IsSome() {
			t.Error("Map should return Some when transforming to float64")
		}
		value, err := result.Unwrap()
		if err != nil {
			t.Errorf("Map result should unwrap without error, got: %v", err)
		}
		if value != 7.5 {
			t.Errorf("Expected 7.5, got %v", value)
		}
	})

	t.Run("does not call transform function when applied to None", func(t *testing.T) {
		m := None[int]()
		called := false
		result := m.Map(func(x int) string {
			called = true
			return "transformed"
		})

		if called {
			t.Error("Map should not call the function for None")
		}
		if !result.IsNone() {
			t.Error("Map should return None when applied to None")
		}
	})
}

func TestMethodFlatMapTypeChange(t *testing.T) {
	t.Run("transforms int to Maybe[string] via method", func(t *testing.T) {
		m := Some(5)
		result := m.FlatMap(func(x int) Maybe[string] {
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

	t.Run("transforms to None when condition fails via method", func(t *testing.T) {
		m := Some(-5)
		result := m.FlatMap(func(x int) Maybe[string] {
			if x > 0 {
				return Some("positive")
			}
			return None[string]()
		})

		if !result.IsNone() {
			t.Error("FlatMap should return None when function returns None")
		}
	})

	t.Run("transforms string to Maybe[int] via method", func(t *testing.T) {
		m := Some("hello")
		result := m.FlatMap(func(s string) Maybe[int] {
			if len(s) > 0 {
				return Some(len(s))
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
		if value != 5 {
			t.Errorf("Expected 5, got %v", value)
		}
	})

	t.Run("does not call transform function when applied to None", func(t *testing.T) {
		m := None[int]()
		called := false
		result := m.FlatMap(func(x int) Maybe[string] {
			called = true
			return Some("transformed")
		})

		if called {
			t.Error("FlatMap should not call the function for None")
		}
		if !result.IsNone() {
			t.Error("FlatMap should return None when applied to None")
		}
	})
}

func TestMethodMatch(t *testing.T) {
	t.Run("calls onSome and not onNone for Some", func(t *testing.T) {
		m := Some(42)
		onNoneCalled := false
		result := m.Match(func(x int) string {
			return fmt.Sprintf("value: %d", x)
		}, func() string {
			onNoneCalled = true
			return "none"
		})

		if onNoneCalled {
			t.Error("onNone should not be called for Some value")
		}
		if result != "value: 42" {
			t.Errorf("Expected 'value: 42', got %v", result)
		}
	})

	t.Run("calls onNone and not onSome for None", func(t *testing.T) {
		m := None[int]()
		onSomeCalled := false
		result := m.Match(func(x int) string {
			onSomeCalled = true
			return fmt.Sprintf("value: %d", x)
		}, func() string {
			return "none"
		})

		if onSomeCalled {
			t.Error("onSome should not be called for None value")
		}
		if result != "none" {
			t.Errorf("Expected 'none', got %v", result)
		}
	})

	t.Run("supports type transformation in return", func(t *testing.T) {
		m := Some(5)
		result := m.Match(func(x int) float64 {
			return float64(x) * 1.5
		}, func() float64 {
			return 0.0
		})

		if result != 7.5 {
			t.Errorf("Expected 7.5, got %v", result)
		}
	})
}

// Helper functions short-circuit on None: the callback must not run.

func TestHelperMapDoesNotCallOnNone(t *testing.T) {
	m := None[int]()
	called := false
	result := Map(m, func(x int) string {
		called = true
		return "transformed"
	})

	if called {
		t.Error("Map helper should not call the function for None")
	}
	if !result.IsNone() {
		t.Error("Map helper should return None when applied to None")
	}
}

func TestHelperFlatMapDoesNotCallOnNone(t *testing.T) {
	m := None[int]()
	called := false
	result := FlatMap(m, func(x int) Maybe[string] {
		called = true
		return Some("transformed")
	})

	if called {
		t.Error("FlatMap helper should not call the function for None")
	}
	if !result.IsNone() {
		t.Error("FlatMap helper should return None when applied to None")
	}
}

func TestHelperMatchBranchExclusivity(t *testing.T) {
	t.Run("does not call onNone for Some value", func(t *testing.T) {
		m := Some(42)
		onNoneCalled := false
		result := Match(m, func(x int) string {
			return fmt.Sprintf("value: %d", x)
		}, func() string {
			onNoneCalled = true
			return "none"
		})

		if onNoneCalled {
			t.Error("onNone should not be called for Some value")
		}
		if result != "value: 42" {
			t.Errorf("Expected 'value: 42', got %v", result)
		}
	})

	t.Run("does not call onSome for None value", func(t *testing.T) {
		m := None[int]()
		onSomeCalled := false
		result := Match(m, func(x int) string {
			onSomeCalled = true
			return fmt.Sprintf("value: %d", x)
		}, func() string {
			return "none"
		})

		if onSomeCalled {
			t.Error("onSome should not be called for None value")
		}
		if result != "none" {
			t.Errorf("Expected 'none', got %v", result)
		}
	})
}

func TestMethodTypeTransformationChain(t *testing.T) {
	t.Run("chains Map and FlatMap across types", func(t *testing.T) {
		result := Some(5).
			Map(func(x int) string { return fmt.Sprintf("%d", x) }).
			FlatMap(func(s string) Maybe[int] {
				if len(s) > 0 {
					return Some(len(s))
				}
				return None[int]()
			}).
			OrElse(0)

		if result != 1 {
			t.Errorf("Expected 1 (length of '5'), got %v", result)
		}
	})

	t.Run("returns default when None flows through chain", func(t *testing.T) {
		result := None[int]().
			Map(func(x int) string { return fmt.Sprintf("%d", x) }).
			FlatMap(func(s string) Maybe[int] { return Some(len(s)) }).
			OrElse(0)

		if result != 0 {
			t.Errorf("Expected default 0, got %v", result)
		}
	})
}
