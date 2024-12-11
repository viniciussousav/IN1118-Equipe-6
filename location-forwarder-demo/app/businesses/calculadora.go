package businesses

type Calculadora struct{}

func (Calculadora) Som(p1, p2 int) int {
	return p1 + p2
}
