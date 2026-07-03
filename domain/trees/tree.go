package trees

import "slices"

type tree struct {
	keyname string
	nodes   []Node
}

func (t *tree) Keyname() string {
	return t.keyname
}

func (t *tree) Nodes() []Node {
	return slices.Clone(t.nodes)
}
