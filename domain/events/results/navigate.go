package results

type navigate struct {
	url string
}

func (n *navigate) URL() string {
	return n.url
}
