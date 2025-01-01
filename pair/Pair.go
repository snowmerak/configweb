package pair

type Pair[A, B any] struct {
	first  A
	second B
}

func New[A, B any](first A, second B) *Pair[A, B] {
	return &Pair[A, B]{first: first, second: second}
}

func (p *Pair[A, B]) First() A {
	return p.first
}

func (p *Pair[A, B]) Second() B {
	return p.second
}
