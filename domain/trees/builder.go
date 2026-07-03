package trees

import (
	"errors"
	"strings"
)

type builder struct {
	keyname string
	nodes   []Node
}

func (b *builder) Create() Builder {
	return &builder{}
}

func (b *builder) WithKeyname(keyname string) Builder {
	b.keyname = strings.TrimSpace(keyname)
	return b
}

func (b *builder) AddNode(node Node) Builder {
	if node != nil {
		b.nodes = append(b.nodes, node)
	}
	return b
}

func (b *builder) WithNodes(nodes []Node) Builder {
	b.nodes = nodes
	return b
}

func (b *builder) Now() (Tree, error) {
	if b.keyname == "" {
		return nil, errors.New("tree keyname is required")
	}
	return &tree{keyname: b.keyname, nodes: b.nodes}, nil
}
