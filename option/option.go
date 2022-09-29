package option

type Option[A any] interface {
	Match(func(), func(A))
}

func Just[A any](a A) Option[A] { return just[A]{a} }

type just[A any] struct{ val A }

func (j just[A]) Match(_ func(), f func(A)) {
	f(j.val)
}

func Nothing[A any]() Option[A] { return nothing[A]{} }

type nothing[A any] struct{}

func (nothing[A]) Match(f func(), _ func(A)) { f() }

func Bind[A, B any](ma Option[A], fn func(A) Option[B]) Option[B] {
	var ret Option[B]
	ma.Match(func() {
		ret = Nothing[B]()
	}, func(a A) {
		ret = fn(a)
	})
	return ret
}

func Map[A, B any](ma Option[A], fn func(A) B) Option[B] {
	return Bind(ma, func(a A) Option[B] {
		return Just(fn(a))
	})
}

func Map2[A, B, C any](ma Option[A], mb Option[B], fn func(A, B) Option[C]) Option[C] {
	return Bind(ma, func(a A) Option[C] {
		return Bind(mb, func(b B) Option[C] {
			return fn(a, b)
		})
	})
}

type Functor[A, B any] func(Option[A]) Option[B]

func With[B, A any](a Option[A], f Functor[A, B]) Option[B] {
	return f(a)
}

func Id[A any](a Option[A]) Option[A] { return a }

func Next[C, A, B any](f func(A) Option[B], g Functor[B, C]) Functor[A, C] {
	return func(a Option[A]) Option[C] {
		return g(Bind(a, f))
	}
}
