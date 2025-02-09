package matrix

import (
	"iter"
	"reflect"
)

type testingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

// Generate generates all possible combinations of values for fields of the type T.
// - len(dims) must equal len(fields(T))
// - dimensions[i].(type) must equal [](fields(T)[i].(type))
// - T must only have exported fields.
func Generate[T any](t testingT, tcase T, dimensions ...any) iter.Seq[T] {
	t.Helper()

	typ := reflect.TypeOf(tcase)

	if typ.Kind() != reflect.Struct {
		t.Fatalf("tcase must be a struct type, got %s", typ.Kind())
	}

	if typ.NumField() == 0 {
		t.Fatalf("tcase must have at least one field")
	}

	if len(dimensions) == 0 {
		t.Fatalf("must supply one dim per tcase field")
	}

	if len(dimensions) != typ.NumField() {
		t.Fatalf("tcase must have same amount of fields as len(dims), got: %d", len(dimensions))
	}

	for i := range typ.NumField() {
		f := typ.Field(i)
		if !f.IsExported() {
			t.Fatalf("tcase must not have unexported fields, got: %s", f.Name)
		}
	}

	rvs := make([][]reflect.Value, len(dimensions))
	for i, dim := range dimensions {
		typ := reflect.TypeOf(dim)
		if typ.Kind() != reflect.Slice {
			t.Fatalf("dims must be slices of tcase field values, got: %s", typ.Kind())
		}
		val := reflect.ValueOf(dim)
		irvs := make([]reflect.Value, 0, val.Len())
		for j := range val.Len() {
			irvs = append(irvs, val.Index(j))
		}
		rvs[i] = irvs
	}

	g := &generator[T]{
		t:    t,
		typ:  typ,
		dims: rvs,
		n:    make([]int, len(rvs)),
	}

	total := len(g.dims[0])
	for _, d := range g.dims[1:] {
		total *= len(d)
	}

	return func(yield func(T) bool) {
		for range total {
			if !yield(g.elem()) {
				return
			}
			g.inc()
		}
	}
}

type generator[T any] struct {
	t    testingT
	typ  reflect.Type
	dims [][]reflect.Value
	n    []int
}

func (g *generator[T]) inc() {
	g.t.Helper()

	for i := 0; i < len(g.n); i++ {
		if g.n[i]+1 < len(g.dims[i]) {
			g.n[i]++
			break
		}
		g.n[i] = 0
	}
}

func (g *generator[T]) elem() T {
	g.t.Helper()

	v := reflect.New(g.typ)

	for i := range v.Elem().NumField() {
		field := v.Elem().Field(i)
		ftype := field.Type()
		val := g.dims[i][g.n[i]]

		if ftype.Kind() == reflect.Pointer {
			ptr := reflect.New(val.Type().Elem())
			ptr.Elem().Set(val.Elem())
			field.Set(ptr)
		} else {
			field.Set(val)
		}
	}
	return v.Elem().Interface().(T)
}
