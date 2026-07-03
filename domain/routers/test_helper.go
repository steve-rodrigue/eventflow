package routers

import (
	"testing"

	"github.com/steve-rodrigue/eventflow/domain/renderables/pages"
	"github.com/steve-rodrigue/eventflow/domain/trees"
)

func mustRequest(t *testing.T, target string, group string, locale string, params ...trees.Param) Request {
	t.Helper()

	request, err := NewBuilder().
		Create().
		WithTargetKeyname(target).
		WithGroupKeyname(group).
		WithLocale(locale).
		WithParams(params).
		Now()

	if err != nil {
		t.Fatalf("expected request, got error: %v", err)
	}

	return request
}

func mustTree(t *testing.T, nodes ...trees.Node) trees.Tree {
	t.Helper()

	tree, err := trees.NewBuilder().
		Create().
		WithKeyname("main").
		WithNodes(nodes).
		Now()

	if err != nil {
		t.Fatalf("expected tree, got error: %v", err)
	}

	return tree
}

func mustNode(t *testing.T, targets ...trees.Target) trees.Node {
	t.Helper()

	node, err := trees.NewNodeBuilder().
		Create().
		WithTargets(targets).
		Now()

	if err != nil {
		t.Fatalf("expected node, got error: %v", err)
	}

	return node
}

func mustTarget(t *testing.T, keyname string, groups ...trees.Group) trees.Target {
	t.Helper()

	target, err := trees.NewTargetBuilder().
		Create().
		WithKeyname(keyname).
		WithGroups(groups).
		Now()

	if err != nil {
		t.Fatalf("expected target, got error: %v", err)
	}

	return target
}

func mustGroup(t *testing.T, keyname string, resources ...trees.Resource) trees.Group {
	t.Helper()

	group, err := trees.NewGroupBuilder().
		Create().
		WithKeyname(keyname).
		WithResources(resources).
		Now()

	if err != nil {
		t.Fatalf("expected group, got error: %v", err)
	}

	return group
}

func mustResource(t *testing.T, locale string, pattern string) trees.Resource {
	t.Helper()

	resource, err := trees.NewResourceBuilder().
		Create().
		WithLocale(locale).
		WithRoute(mustRoute(t, pattern)).
		WithPage(pages.MustPage(t)).
		Now()

	if err != nil {
		t.Fatalf("expected resource, got error: %v", err)
	}

	return resource
}

func mustResourceWithChildren(t *testing.T, locale string, pattern string, children trees.Tree) trees.Resource {
	t.Helper()

	resource, err := trees.NewResourceBuilder().
		Create().
		WithLocale(locale).
		WithRoute(mustRoute(t, pattern)).
		WithPage(pages.MustPage(t)).
		WithChildren(children).
		Now()

	if err != nil {
		t.Fatalf("expected resource, got error: %v", err)
	}

	return resource
}

func mustRoute(t *testing.T, pattern string) trees.Route {
	t.Helper()

	route, err := trees.NewRouteBuilder().
		Create().
		WithPattern(pattern).
		Now()

	if err != nil {
		t.Fatalf("expected route, got error: %v", err)
	}

	return route
}

func mustParam(t *testing.T, keyname string, value string) trees.Param {
	t.Helper()

	param, err := trees.NewParamBuilder().
		Create().
		WithKeyname(keyname).
		WithValue(value).
		Now()

	if err != nil {
		t.Fatalf("expected param, got error: %v", err)
	}

	return param
}
