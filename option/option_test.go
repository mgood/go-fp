package option_test

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/mgood/go-fp/option"
)

func TestLeftIdentity(t *testing.T) {
	// unit(x) >>= f ↔ f(x)
	f := func(x int) option.Option[int] {
		return option.Just(x * 2)
	}
	if err := quick.CheckEqual(f, func(x int) option.Option[int] {
		unit := option.Just(x)
		return option.Bind(unit, f)
	}, nil); err != nil {
		t.Errorf("%v", err)
	}
}

func TestRightIdentity(t *testing.T) {
	// ma >>= unit ↔ ma
	if err := quick.CheckEqual(option.Just[int], func(x int) option.Option[int] {
		ma := option.Just(x)
		return option.Bind(ma, option.Just[int])
	}, nil); err != nil {
		t.Errorf("%v", err)
	}
}

func TestChain(t *testing.T) {
	r := option.With(
		option.Just(1),
		option.Next(func(x int) option.Option[float32] {
			return option.Just(float32(x * 2))
		}, option.Next(func(x float32) option.Option[float64] {
			return option.Just(float64(x))
		}, option.Next(func(x float64) option.Option[string] {
			return option.Just(fmt.Sprintf("%.1f", x))
		}, option.Next(func(x string) option.Option[string] {
			return option.Just(x + x)
		}, option.Id[string],
		)))))
	r.Match(func() {
		t.Errorf("result should not be Nothing")
	}, func(s string) {
		expected := "2.02.0"
		if s != expected {
			t.Errorf("result should be %v, but got %v", expected, s)
		}
	})
}

func TestAdd(t *testing.T) {
	add := func(ma, mb option.Option[int]) option.Option[int] {
		return option.Map2(ma, mb, func(a, b int) option.Option[int] {
			return option.Just(a + b)
		})
	}
	if err := quick.CheckEqual(func(a, b int) option.Option[int] {
		return option.Just(a + b)
	}, func(a, b int) option.Option[int] {
		return add(option.Just(a), option.Just(b))
	}, nil); err != nil {
		t.Errorf("%v", err)
	}
}
