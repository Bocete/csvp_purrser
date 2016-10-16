package func_purrser

type Node interface {
}

type Expr interface {
	Node
}

type TabRange struct {
	descriptor string
}

type FunctionCall struct {
	name string
	args []Expr
}
