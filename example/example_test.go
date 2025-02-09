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
