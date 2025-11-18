# Go Maybe

A generic Maybe/Option type implementation for Go, providing a safe way to handle potentially absent values without using pointers or nil checks.

## Features

- **Type-safe**: Uses Go generics to provide compile-time type safety
- **Zero dependencies**: Pure Go implementation with no external dependencies
- **Comprehensive API**: Includes `Map`, `FlatMap`, `Filter`, and more
- **Type transformations**: Helper functions support transforming between different types
- **Safe unwrapping**: Multiple ways to extract values with proper error handling

## Installation

```bash
go get github.com/zodimo/go-maybe
```

## Usage

### Creating Maybe Values

```go
import "github.com/zodimo/go-maybe"

// Create a Some value
value := maybe.Some(42)

// Create a None value
empty := maybe.None[int]()

// Create a new empty Maybe (same as None)
m := maybe.NewMaybe[string]()
```

### Checking Values

```go
m := maybe.Some(42)

if m.IsSome() {
    // Value is present
}

if m.IsNone() {
    // Value is absent
}
```

### Extracting Values

```go
// Safe unwrapping with error
value, err := m.Unwrap()
if err != nil {
    // Handle error
}

// Unsafe unwrapping (panics if None)
value := m.UnwrapUnsafe() // Use with caution!

// UnwrapOr: provide a default value if None
m := maybe.None[int]()
value := m.UnwrapOr(100) // Returns 100
```

### Transforming Values

#### Methods (Same Type)

The methods `Map`, `FlatMap`, and `Filter` work on the same type:

```go
// Map: transform the value if present (same type)
m := maybe.Some(5)
doubled := m.Map(func(x int) int { return x * 2 })
// doubled is Some(10)

// FlatMap: transform to another Maybe (same type)
m := maybe.Some(5)
result := m.FlatMap(func(x int) maybe.Maybe[int] {
    if x > 0 {
        return maybe.Some(x * 2)
    }
    return maybe.None[int]()
})

// Filter: keep value only if predicate is true
m := maybe.Some(10)
filtered := m.Filter(func(x int) bool { return x > 5 })
// filtered is Some(10)
```

#### Helper Functions (Type Transformations)

The standalone helper functions `Map`, `FlatMap`, and `Filter` support transforming between different types:

```go
// Map: transform int to string
m := maybe.Some(42)
result := maybe.Map(m, func(x int) string {
    return fmt.Sprintf("number: %d", x)
})
// result is Maybe[string] with value "number: 42"

// Map: transform string to int
m := maybe.Some("hello")
result := maybe.Map(m, func(s string) int {
    return len(s)
})
// result is Maybe[int] with value 5

// FlatMap: transform int to Maybe[string]
m := maybe.Some(5)
result := maybe.FlatMap(m, func(x int) maybe.Maybe[string] {
    if x > 0 {
        return maybe.Some("positive")
    }
    return maybe.None[string]()
})
// result is Maybe[string] with value "positive"
```

### Providing Default Values

```go
// OrElse: provide a default value
m := maybe.None[int]()
value := m.OrElse(100) // Returns 100

// OrElseGet: compute default value lazily
m := maybe.None[int]()
value := m.OrElseGet(func() int {
    return expensiveComputation()
})

// OrElseError: return an error if None
m := maybe.None[int]()
value, err := m.OrElseError(errors.New("value not found"))
```

### Chaining Operations

```go
result := maybe.Some(5).
    Map(func(x int) int { return x * 2 }).
    Map(func(x int) int { return x + 1 }).
    Filter(func(x int) bool { return x > 10 }).
    OrElse(0)
```

## API Reference

### Types

- `Maybe[T]`: A generic type that represents either a value (`Some`) or no value (`None`)

### Functions

- `Some[T any](value T) Maybe[T]`: Creates a `Maybe` with a value
- `None[T any]() Maybe[T]`: Creates an empty `Maybe`
- `NewMaybe[T any]() Maybe[T]`: Creates a new empty `Maybe` (alias for `None`)
- `Map[T any, R any](m Maybe[T], f func(T) R) Maybe[R]`: Transforms a `Maybe[T]` to `Maybe[R]` by applying function `f` if the value is present
- `FlatMap[T any, R any](m Maybe[T], f func(T) Maybe[R]) Maybe[R]`: Transforms a `Maybe[T]` to `Maybe[R]` by applying function `f` that returns a `Maybe[R]` if the value is present

### Methods

- `IsSome() bool`: Returns `true` if the `Maybe` contains a value
- `IsNone() bool`: Returns `true` if the `Maybe` is empty
- `Unwrap() (T, error)`: Returns the value and an error (error is non-nil if `None`)
- `UnwrapUnsafe() T`: Returns the value, panics if `None`
- `UnwrapOr(defaultValue T) T`: Returns the value if present, otherwise returns `defaultValue`
- `Map(f func(T) T) Maybe[T]`: Transforms the value if present, returns `None` otherwise
- `FlatMap(f func(T) Maybe[T]) Maybe[T]`: Transforms to another `Maybe` if present
- `Filter(f func(T) bool) Maybe[T]`: Keeps the value only if the predicate returns `true`
- `OrElse(elseValue T) T`: Returns the value if present, otherwise returns `elseValue`
- `OrElseGet(f func() T) T`: Returns the value if present, otherwise calls `f()` and returns its result
- `OrElseError(err error) (T, error)`: Returns the value and `nil` error if present, otherwise returns zero value and `err`

## Examples

### Handling Optional Return Values

```go
func divide(a, b int) maybe.Maybe[int] {
    if b == 0 {
        return maybe.None[int]()
    }
    return maybe.Some(a / b)
}

result := divide(10, 2)
value, err := result.Unwrap()
if err != nil {
    log.Println("Division failed")
} else {
    log.Printf("Result: %d", value)
}
```

### Safe Map Access

```go
func getValue(m map[string]int, key string) maybe.Maybe[int] {
    if val, ok := m[key]; ok {
        return maybe.Some(val)
    }
    return maybe.None[int]()
}

m := map[string]int{"foo": 42}
value := getValue(m, "foo").OrElse(0)
```

### Type Transformations with Helper Functions

```go
// Transform a number to a formatted string
num := maybe.Some(42)
formatted := maybe.Map(num, func(x int) string {
    return fmt.Sprintf("Value: %d", x)
})
// formatted is Maybe[string] with value "Value: 42"

// Parse a string to an integer with validation
str := maybe.Some("123")
parsed := maybe.FlatMap(str, func(s string) maybe.Maybe[int] {
    if val, err := strconv.Atoi(s); err == nil {
        return maybe.Some(val)
    }
    return maybe.None[int]()
})
// parsed is Maybe[int] with value 123
```

### Chaining Transformations

```go
// Using methods (same type)
result := maybe.Some("hello").
    Map(func(s string) string { return strings.ToUpper(s) }).
    Map(func(s string) string { return s + " WORLD" }).
    Filter(func(s string) bool { return len(s) > 5 }).
    OrElse("default")

// Using helper functions (type transformations)
m := maybe.Some(42)
strMaybe := maybe.Map(m, func(x int) string { return fmt.Sprintf("%d", x) })
intMaybe := maybe.FlatMap(strMaybe, func(s string) maybe.Maybe[int] {
    if len(s) > 0 {
        return maybe.Some(len(s))
    }
    return maybe.None[int]()
})
result := intMaybe.OrElse(0)
```

## License

See LICENSE file for details.
