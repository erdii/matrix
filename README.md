# Matrix

A small matrix test generator implementation for go.

## Why

The general approach to writing parametric tests in go is often a table-driven test, that spells out all variations (the cartesian product) of the table and looks somewhat like this:

```golang
func TestExample(t *testing.T) {
  type tcase struct {
    X, Y bool
  }

  for _, tcase := range []tcase{
    {false, false},
    {false, true},
    {true, false},
    {true, true},
  } {
    t.Run(fmt.Sprintf("%t_%t", tcase.X, tcase.Y), func(t *testing.T) {
      // Imaginary test code.
    })
  }
}
```

This can become very tedious, especially when more dimensions get added to the test matrix.

## What

In this package I propose a solution that works like this:
- Define a struct with one field per test dimension. (Just like the example above).
- Call `matrix.Generate` passing your struct type and slices of values per dimension.
- `matrix.Generate` returns an iterator that you can `range` over.

## Example

It's easier to explain [through code](./example/example_test.go), so here we go:

```go
package example

import (
  "fmt"
  "testing"

  "github.com/erdii/matrix"
)

func TestExample(t *testing.T) {
  t.Parallel()

  // Imagine that you want to test all possible combinations
  // of Priyanka, Pedro and Paul
  // eating Pizza, Spaghetti, Bananas and Sushi:
  type tcase struct {
    Who  string
    Food string
  }

  for tcase := range matrix.Generate(t, tcase{},
    []string{"Priyanka", "Pedro", "Paul"},
    []string{"Pizza", "Spaghetti", "Bananas", "Sushi"}) {
    t.Run(fmt.Sprintf("%s_%s", tcase.Who, tcase.Food), func(t *testing.T) {
      // Imaginary test code.
    })
  }
}
```

Running the example test with `go test -v ./example` will print:

```
=== RUN   TestExample
=== PAUSE TestExample
=== CONT  TestExample
=== RUN   TestExample/Priyanka_Pizza
=== RUN   TestExample/Pedro_Pizza
=== RUN   TestExample/Paul_Pizza
=== RUN   TestExample/Priyanka_Spaghetti
=== RUN   TestExample/Pedro_Spaghetti
=== RUN   TestExample/Paul_Spaghetti
=== RUN   TestExample/Priyanka_Bananas
=== RUN   TestExample/Pedro_Bananas
=== RUN   TestExample/Paul_Bananas
=== RUN   TestExample/Priyanka_Sushi
=== RUN   TestExample/Pedro_Sushi
=== RUN   TestExample/Paul_Sushi
--- PASS: TestExample (0.00s)
    --- PASS: TestExample/Priyanka_Pizza (0.00s)
    --- PASS: TestExample/Pedro_Pizza (0.00s)
    --- PASS: TestExample/Paul_Pizza (0.00s)
    --- PASS: TestExample/Priyanka_Spaghetti (0.00s)
    --- PASS: TestExample/Pedro_Spaghetti (0.00s)
    --- PASS: TestExample/Paul_Spaghetti (0.00s)
    --- PASS: TestExample/Priyanka_Bananas (0.00s)
    --- PASS: TestExample/Pedro_Bananas (0.00s)
    --- PASS: TestExample/Paul_Bananas (0.00s)
    --- PASS: TestExample/Priyanka_Sushi (0.00s)
    --- PASS: TestExample/Pedro_Sushi (0.00s)
    --- PASS: TestExample/Paul_Sushi (0.00s)
PASS
ok  	github.com/erdii/matrix/example	0.001s
```

---

Coincidentally this package is designed in a similar way to rust's [test_matrix](https://docs.rs/test-case/latest/test_case/#test-matrix).
