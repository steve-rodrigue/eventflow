package trees

import "slices"

type group struct {
	keyname   string
	resources []Resource
}

func (g *group) Keyname() string {
	return g.keyname
}

func (g *group) HasResources() bool {
	return len(g.resources) > 0
}

func (g *group) Resources() []Resource {
	return slices.Clone(g.resources)
}

func (g *group) Resource(locale string) (Resource, bool) {
	for _, resource := range g.resources {
		if resource.Locale() == locale {
			return resource, true
		}
	}
	return nil, false
}
