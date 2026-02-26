package maybe

import "testing"

func TestFromPtrDereferenced(t *testing.T) {
	t.Run("returns Some when pointer is non-nil int", func(t *testing.T) {
		val := 42
		m := FromPtrDereferenced(&val)

		if !m.IsSome() {
			t.Error("FromPtrDereferenced with non-nil pointer should return Some")
		}
		if m.IsNone() {
			t.Error("FromPtrDereferenced with non-nil pointer should not return None")
		}

		got, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got != 42 {
			t.Errorf("Expected 42, got %v", got)
		}
	})

	t.Run("returns Some when pointer is non-nil string", func(t *testing.T) {
		val := "hello"
		m := FromPtrDereferenced(&val)

		if !m.IsSome() {
			t.Error("FromPtrDereferenced with non-nil pointer should return Some")
		}

		got, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got != "hello" {
			t.Errorf("Expected 'hello', got %v", got)
		}
	})

	t.Run("returns Some with zero value pointer", func(t *testing.T) {
		val := 0
		m := FromPtrDereferenced(&val)

		if !m.IsSome() {
			t.Error("FromPtrDereferenced with pointer to zero value should return Some")
		}

		got, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got != 0 {
			t.Errorf("Expected 0, got %v", got)
		}
	})

	t.Run("returns Some with empty string pointer", func(t *testing.T) {
		val := ""
		m := FromPtrDereferenced(&val)

		if !m.IsSome() {
			t.Error("FromPtrDereferenced with pointer to empty string should return Some")
		}

		got, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got != "" {
			t.Errorf("Expected empty string, got %v", got)
		}
	})

	t.Run("returns None when pointer is nil int", func(t *testing.T) {
		var ptr *int
		m := FromPtrDereferenced(ptr)

		if m.IsSome() {
			t.Error("FromPtrDereferenced with nil pointer should return None")
		}
		if !m.IsNone() {
			t.Error("FromPtrDereferenced with nil pointer should return None")
		}
	})

	t.Run("returns None when pointer is nil string", func(t *testing.T) {
		var ptr *string
		m := FromPtrDereferenced(ptr)

		if m.IsSome() {
			t.Error("FromPtrDereferenced with nil pointer should return None")
		}
		if !m.IsNone() {
			t.Error("FromPtrDereferenced with nil pointer should return None")
		}
	})

	t.Run("value is a copy not a reference", func(t *testing.T) {
		val := 42
		m := FromPtrDereferenced(&val)
		val = 100 // mutate original

		got, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got != 42 {
			t.Errorf("Expected 42 (original value), got %v — value should be copied, not referenced", got)
		}
	})

	t.Run("works with struct type", func(t *testing.T) {
		type Point struct {
			X, Y int
		}
		val := Point{X: 1, Y: 2}
		m := FromPtrDereferenced(&val)

		if !m.IsSome() {
			t.Error("FromPtrDereferenced with non-nil struct pointer should return Some")
		}

		got, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got.X != 1 || got.Y != 2 {
			t.Errorf("Expected Point{1, 2}, got %v", got)
		}
	})
}

func TestToPtr(t *testing.T) {
	t.Run("returns non-nil pointer for Some int", func(t *testing.T) {
		m := Some(42)
		ptr := m.ToPtr()

		if ptr == nil {
			t.Fatal("ToPtr should return non-nil pointer for Some")
		}
		if *ptr != 42 {
			t.Errorf("Expected *ptr to be 42, got %v", *ptr)
		}
	})

	t.Run("returns non-nil pointer for Some string", func(t *testing.T) {
		m := Some("hello")
		ptr := m.ToPtr()

		if ptr == nil {
			t.Fatal("ToPtr should return non-nil pointer for Some")
		}
		if *ptr != "hello" {
			t.Errorf("Expected *ptr to be 'hello', got %v", *ptr)
		}
	})

	t.Run("returns non-nil pointer for Some zero value", func(t *testing.T) {
		m := Some(0)
		ptr := m.ToPtr()

		if ptr == nil {
			t.Fatal("ToPtr should return non-nil pointer for Some with zero value")
		}
		if *ptr != 0 {
			t.Errorf("Expected *ptr to be 0, got %v", *ptr)
		}
	})

	t.Run("returns non-nil pointer for Some empty string", func(t *testing.T) {
		m := Some("")
		ptr := m.ToPtr()

		if ptr == nil {
			t.Fatal("ToPtr should return non-nil pointer for Some with empty string")
		}
		if *ptr != "" {
			t.Errorf("Expected *ptr to be empty string, got %v", *ptr)
		}
	})

	t.Run("returns nil pointer for None int", func(t *testing.T) {
		m := None[int]()
		ptr := m.ToPtr()

		if ptr != nil {
			t.Errorf("ToPtr should return nil for None, got %v", *ptr)
		}
	})

	t.Run("returns nil pointer for None string", func(t *testing.T) {
		m := None[string]()
		ptr := m.ToPtr()

		if ptr != nil {
			t.Errorf("ToPtr should return nil for None, got %v", *ptr)
		}
	})

	t.Run("returned pointer is independent of Maybe", func(t *testing.T) {
		m := Some(42)
		ptr := m.ToPtr()
		*ptr = 100 // mutate via pointer

		// The Maybe value should be unchanged
		got, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got != 42 {
			t.Errorf("Expected Maybe value to remain 42, got %v — pointer should be independent", got)
		}
	})

	t.Run("works with struct type", func(t *testing.T) {
		type Point struct {
			X, Y int
		}
		m := Some(Point{X: 1, Y: 2})
		ptr := m.ToPtr()

		if ptr == nil {
			t.Fatal("ToPtr should return non-nil pointer for Some struct")
		}
		if ptr.X != 1 || ptr.Y != 2 {
			t.Errorf("Expected &Point{1, 2}, got %v", *ptr)
		}
	})
}

func TestFromPtrDereferencedToPtr_Roundtrip(t *testing.T) {
	t.Run("Some roundtrip preserves value", func(t *testing.T) {
		val := 42
		m := FromPtrDereferenced(&val)
		ptr := m.ToPtr()

		if ptr == nil {
			t.Fatal("Roundtrip should produce non-nil pointer")
		}
		if *ptr != 42 {
			t.Errorf("Roundtrip expected 42, got %v", *ptr)
		}
	})

	t.Run("None roundtrip preserves nil", func(t *testing.T) {
		var ptr *int
		m := FromPtrDereferenced(ptr)
		result := m.ToPtr()

		if result != nil {
			t.Errorf("Roundtrip of nil should produce nil, got %v", *result)
		}
	})

	t.Run("ToPtr then FromPtrDereferenced roundtrip", func(t *testing.T) {
		m := Some("test")
		ptr := m.ToPtr()
		m2 := FromPtrDereferenced(ptr)

		if !m2.IsSome() {
			t.Error("Roundtrip should produce Some")
		}

		got, err := m2.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got != "test" {
			t.Errorf("Expected 'test', got %v", got)
		}
	})

	t.Run("None ToPtr then FromPtrDereferenced roundtrip", func(t *testing.T) {
		m := None[int]()
		ptr := m.ToPtr()
		m2 := FromPtrDereferenced(ptr)

		if !m2.IsNone() {
			t.Error("Roundtrip of None should produce None")
		}
	})
}

func TestFromPtr(t *testing.T) {
	t.Run("returns Some when pointer is non-nil int", func(t *testing.T) {
		val := 42
		m := FromPtr(&val)

		if !m.IsSome() {
			t.Error("FromPtr with non-nil pointer should return Some")
		}

		got, err := m.Unwrap()
		if err != nil {
			t.Errorf("Unwrap should not return error, got: %v", err)
		}
		if got != &val {
			t.Errorf("Expected pointer to val, got %v", got)
		}
	})

	t.Run("returns None when pointer is nil string", func(t *testing.T) {
		var ptr *string
		m := FromPtr(ptr)

		if m.IsSome() {
			t.Error("FromPtr with nil pointer should return None")
		}
		if !m.IsNone() {
			t.Error("FromPtr with nil pointer should return None")
		}
	})
}

func TestFromInterface(t *testing.T) {
	t.Run("returns Some when interface has concrete value", func(t *testing.T) {
		type customStruct struct{}
		var val any = &customStruct{}
		m := FromInterface(val)

		if !m.IsSome() {
			t.Error("FromInterface with concrete value should return Some")
		}

		got, unwrapErr := m.Unwrap()
		if unwrapErr != nil {
			t.Errorf("Unwrap should not return error, got: %v", unwrapErr)
		}
		if got != val {
			t.Errorf("Expected specific value, got %v", got)
		}
	})

	t.Run("returns None when interface is untyped nil", func(t *testing.T) {
		var err error = nil
		m := FromInterface(err)

		if m.IsSome() {
			t.Error("FromInterface with nil should return None")
		}
		if !m.IsNone() {
			t.Error("FromInterface with nil should return None")
		}
	})

	t.Run("returns Some with typed string", func(t *testing.T) {
		var val any = "hello"
		m := FromInterface(val)

		if !m.IsSome() {
			t.Error("FromInterface with string should return Some")
		}

		got, unwrapErr := m.Unwrap()
		if unwrapErr != nil {
			t.Errorf("Unwrap should not return error, got: %v", unwrapErr)
		}
		if got != "hello" {
			t.Errorf("Expected specific string, got %v", got)
		}
	})

	t.Run("returns None when custom interface is nil", func(t *testing.T) {
		type Person interface {
			Name() string
		}
		var person Person = nil
		m := FromInterface(person)

		if m.IsSome() {
			t.Error("FromInterface with nil custom interface should return None")
		}
		if !m.IsNone() {
			t.Error("FromInterface with nil custom interface should return None")
		}
	})
}
