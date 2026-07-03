package trees

type param struct {
	keyname string
	value   string
}

func (p *param) Keyname() string {
	return p.keyname
}

func (p *param) Value() string {
	return p.value
}
