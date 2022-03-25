package o

type Option[T any] func(*T)

func Apply[T any](t *T, options ...Option[T]) {
	for _, apply := range options {
		apply(t)
	}
}

type OptionGroup[T any] []Option[T]

func Options[T any](options ...Option[T]) *OptionGroup[T] {
	g := OptionGroup[T](options)
	return &g
}

func (o *OptionGroup[T]) Add(option Option[T], options ...Option[T]) *OptionGroup[T] {
	*o = append(*o, options...)
	return o
}

func (o *OptionGroup[T]) AddC(ok bool, option Option[T], options ...Option[T]) *OptionGroup[T] {
	if ok {
		*o = append(*o, options...)
	}
	return o
}

func (o *OptionGroup[T]) O() []Option[T] {
	return *o
}
