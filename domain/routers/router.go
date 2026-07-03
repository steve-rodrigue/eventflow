package routers

import (
	"errors"
	"strings"

	"github.com/steve-rodrigue/eventflow/domain/trees"
)

type router struct{}

func (r *router) URI(tree trees.Tree, request Request) (string, error) {
	if tree == nil {
		return "", errors.New("tree is required")
	}

	if request == nil {
		return "", errors.New("route request is required")
	}

	resource, ok := r.findResource(tree, request)
	if !ok {
		return "", errors.New("resource not found")
	}

	return r.renderPattern(resource.Route().Pattern(), request.Params()), nil
}

func (r *router) findResource(tree trees.Tree, request Request) (trees.Resource, bool) {
	for _, node := range tree.Nodes() {
		for _, target := range node.Targets() {
			if target.Keyname() != request.TargetKeyname() {
				continue
			}

			group, ok := target.Group(request.GroupKeyname())
			if !ok {
				continue
			}

			resource, ok := group.Resource(request.Locale())
			if !ok {
				continue
			}

			if resource.HasChildren() {
				found, ok := r.findResource(resource.Children(), request)
				if ok {
					return found, true
				}
			}

			return resource, true
		}
	}

	return nil, false
}

func (r *router) renderPattern(pattern string, params []trees.Param) string {
	uri := pattern

	for _, param := range params {
		uri = strings.ReplaceAll(
			uri,
			"{"+param.Keyname()+"}",
			param.Value(),
		)
	}

	return uri
}
