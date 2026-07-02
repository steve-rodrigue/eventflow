package components

import (
	"testing"

	"github.com/steve-rodrigue/eventflow/domain/events"
	"github.com/steve-rodrigue/eventflow/domain/renderables"
	"github.com/steve-rodrigue/eventflow/domain/renderables/pages/templates"
)

func TestRendererRenderSimpleComponent(t *testing.T) {
	component := mustComponent(t, "message", mustTemplate(t, "message", "<p>{text}</p>"))

	result := NewRenderer(templates.NewMustacheRenderer()).Render(component, renderables.Params{
		"text": "Hello",
	})

	expected := "<p>Hello</p>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildWithRenderableParams(t *testing.T) {
	child := mustComponent(t, "user_card", mustTemplate(t, "user_card", "<span>{name}</span>"))
	parent := mustComponent(t, "page", mustTemplate(t, "page", "<main>{title}{user_card}</main>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).Render(parent, renderables.Params{
		"title": "User: ",
		"user_card": renderables.Params{
			"name": "Steve",
		},
	})

	expected := "<main>User: <span>Steve</span></main>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildWithMapParams(t *testing.T) {
	child := mustComponent(t, "user_card", mustTemplate(t, "user_card", "<span>{name}</span>"))
	parent := mustComponent(t, "page", mustTemplate(t, "page", "<main>{user_card}</main>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).Render(parent, renderables.Params{
		"user_card": map[string]any{
			"name": "Steve",
		},
	})

	expected := "<main><span>Steve</span></main>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildWithRenderableParamsList(t *testing.T) {
	child := mustComponent(t, "item", mustTemplate(t, "item", "<li>{label}</li>"))
	parent := mustComponent(t, "list", mustTemplate(t, "list", "<ul>{item}</ul>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).Render(parent, renderables.Params{
		"item": []renderables.Params{
			{"label": "One"},
			{"label": "Two"},
		},
	})

	expected := "<ul><li>One</li><li>Two</li></ul>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildWithMapParamsList(t *testing.T) {
	child := mustComponent(t, "item", mustTemplate(t, "item", "<li>{label}</li>"))
	parent := mustComponent(t, "list", mustTemplate(t, "list", "<ul>{item}</ul>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).Render(parent, renderables.Params{
		"item": []map[string]any{
			{"label": "One"},
			{"label": "Two"},
		},
	})

	expected := "<ul><li>One</li><li>Two</li></ul>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildWithDefaultValue(t *testing.T) {
	child := mustComponent(t, "value", mustTemplate(t, "value", "<span>{text}</span>"))
	parent := mustComponent(t, "page", mustTemplate(t, "page", "<main>{value}</main>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).Render(parent, renderables.Params{
		"value": "raw value",
	})

	expected := "<main>raw value</main>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderSkipsChildWhenNoParamsExist(t *testing.T) {
	child := mustComponent(t, "user_card", mustTemplate(t, "user_card", "<span>{name}</span>"))
	parent := mustComponent(t, "page", mustTemplate(t, "page", "<main>{title}{user_card}</main>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).Render(parent, renderables.Params{
		"title": "Dashboard",
	})

	expected := "<main>Dashboard</main>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderList(t *testing.T) {
	first := mustComponent(t, "first", mustTemplate(t, "first", "<p>{value}</p>"))
	second := mustComponent(t, "second", mustTemplate(t, "second", "<span>{value}</span>"))

	result := NewRenderer(templates.NewMustacheRenderer()).RenderList(
		[]Component{first, second},
		[]renderables.Params{
			{"value": "One"},
			{"value": "Two"},
		},
	)

	expected := "<p>One</p><span>Two</span>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderListUsesEmptyParamsWhenMissing(t *testing.T) {
	first := mustComponent(t, "first", mustTemplate(t, "first", "<p>{value}</p>"))
	second := mustComponent(t, "second", mustTemplate(t, "second", "<span>{value}</span>"))

	result := NewRenderer(templates.NewMustacheRenderer()).RenderList(
		[]Component{first, second},
		[]renderables.Params{
			{"value": "One"},
		},
	)

	expected := "<p>One</p><span></span>"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestBuilderCreatePreservesStylableBuilder(t *testing.T) {
	builder := NewBuilder(renderables.NewStylableBuilder())
	created := builder.Create()

	component, err := created.
		WithKeyname("test").
		WithTemplate(mustTemplate(t, "test", "<p>Hello</p>")).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if component.Keyname() != "test" {
		t.Fatalf("expected keyname test, got %s", component.Keyname())
	}
}

func TestBuilderNowWithStyleEventsAndChildren(t *testing.T) {
	style := mustTemplate(t, "style", ".user { color: {color}; }")
	child := mustComponent(t, "child", mustTemplate(t, "child", "<span>{name}</span>"))
	event := mustEvent(t, "clicked")

	component, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("  user_card  ").
		WithTemplate(mustTemplate(t, "user_card", "<div>{child}</div>")).
		WithStyle(style).
		AddEvent(event).
		AddEvent(nil).
		AddChild(child).
		AddChild(nil).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if component.Keyname() != "user_card" {
		t.Fatalf("expected trimmed keyname user_card, got %q", component.Keyname())
	}

	if !component.HasStyle() {
		t.Fatal("expected component to have style")
	}

	if component.Style() != style {
		t.Fatal("expected style to match")
	}

	if !component.HasEvents() {
		t.Fatal("expected component to have events")
	}

	if len(component.Events()) != 1 {
		t.Fatalf("expected 1 event, got %d", len(component.Events()))
	}

	if !component.HasChildren() {
		t.Fatal("expected component to have children")
	}

	if len(component.Children()) != 1 {
		t.Fatalf("expected 1 child, got %d", len(component.Children()))
	}
}

func TestBuilderWithEventsAndWithChildren(t *testing.T) {
	event := mustEvent(t, "clicked")
	child := mustComponent(t, "child", mustTemplate(t, "child", "<span>Child</span>"))

	component, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("parent").
		WithTemplate(mustTemplate(t, "parent", "<div>{child}</div>")).
		WithEvents([]events.Event{event}).
		WithChildren([]Component{child}).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(component.Events()) != 1 {
		t.Fatalf("expected 1 event, got %d", len(component.Events()))
	}

	if len(component.Children()) != 1 {
		t.Fatalf("expected 1 child, got %d", len(component.Children()))
	}
}

func TestBuilderNowReturnsErrorWhenKeynameIsMissing(t *testing.T) {
	_, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("   ").
		WithTemplate(mustTemplate(t, "test", "<p>Hello</p>")).
		Now()

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "component keyname is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuilderNowReturnsErrorWhenStylableBuilderIsMissing(t *testing.T) {
	_, err := NewBuilder(nil).
		Create().
		WithKeyname("test").
		WithTemplate(mustTemplate(t, "test", "<p>Hello</p>")).
		Now()

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "component stylable builder is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuilderNowReturnsStylableBuilderError(t *testing.T) {
	_, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("test").
		Now()

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "renderable template is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestComponentEventsReturnsClone(t *testing.T) {
	first := mustEvent(t, "first")
	second := mustEvent(t, "second")

	component, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("test").
		WithTemplate(mustTemplate(t, "test", "<p>Hello</p>")).
		WithEvents([]events.Event{first}).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cloned := component.Events()
	cloned[0] = second

	if component.Events()[0].Keyname() != "first" {
		t.Fatal("expected events to be cloned")
	}
}

func TestComponentChildrenReturnsClone(t *testing.T) {
	first := mustComponent(t, "first", mustTemplate(t, "first", "<p>First</p>"))
	second := mustComponent(t, "second", mustTemplate(t, "second", "<p>Second</p>"))

	component, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("test").
		WithTemplate(mustTemplate(t, "test", "<div>{first}</div>")).
		WithChildren([]Component{first}).
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cloned := component.Children()
	cloned[0] = second

	if component.Children()[0].Keyname() != "first" {
		t.Fatal("expected children to be cloned")
	}
}

func TestRendererRenderStyleWithoutStyle(t *testing.T) {
	component := mustComponent(t, "message", mustTemplate(t, "message", "<p>{text}</p>"))

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyle(component, renderables.Params{
		"text": "Hello",
	})

	if result != "" {
		t.Fatalf("expected empty style, got %q", result)
	}
}

func TestRendererRenderStyle(t *testing.T) {
	component, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("myComponent").
		WithTemplate(mustTemplate(t, "template", "<div>Hello</div>")).
		WithStyle(mustTemplate(t, "style", ".myComponent { background-color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected component, got error: %v", err)
	}

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyle(component, renderables.Params{
		"color": "#000000",
	})

	expected := ".myComponent { background-color: #000000; }"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderStyleScopesChildStyle(t *testing.T) {
	child, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("myComponent").
		WithTemplate(mustTemplate(t, "child_template", "<div>Child</div>")).
		WithStyle(mustTemplate(t, "child_style", ".myComponent { background-color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected child, got error: %v", err)
	}

	parent := mustComponent(t, "myParent", mustTemplate(t, "parent_template", "<section>{myComponent}</section>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyle(parent, renderables.Params{
		"myComponent": renderables.Params{
			"color": "#000000",
		},
	})

	expected := ".myParent .myComponent { background-color: #000000; }"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderStyleScopesNestedChildStyle(t *testing.T) {
	child, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("myComponent").
		WithTemplate(mustTemplate(t, "child_template", "<div>Child</div>")).
		WithStyle(mustTemplate(t, "child_style", ".myComponent { background-color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected child, got error: %v", err)
	}

	parent := mustComponent(t, "myParent", mustTemplate(t, "parent_template", "<section>{myComponent}</section>"), child)

	grandParent := mustComponent(t, "myGrandParent", mustTemplate(t, "grand_template", "<main>{myParent}</main>"), parent)

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyle(grandParent, renderables.Params{
		"myParent": renderables.Params{
			"myComponent": renderables.Params{
				"color": "#000000",
			},
		},
	})

	expected := ".myGrandParent .myParent .myComponent { background-color: #000000; }"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderStyleList(t *testing.T) {
	first, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("first").
		WithTemplate(mustTemplate(t, "first_template", "<div>First</div>")).
		WithStyle(mustTemplate(t, "first_style", ".first { color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected first component, got error: %v", err)
	}

	second, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("second").
		WithTemplate(mustTemplate(t, "second_template", "<div>Second</div>")).
		WithStyle(mustTemplate(t, "second_style", ".second { color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected second component, got error: %v", err)
	}

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyleList(
		[]Component{first, second},
		[]renderables.Params{
			{"color": "red"},
			{"color": "blue"},
		},
	)

	expected := ".first { color: red; }\n.second { color: blue; }"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderStyleListSkipsEmptyStyles(t *testing.T) {
	first := mustComponent(t, "first", mustTemplate(t, "first_template", "<div>First</div>"))

	second, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("second").
		WithTemplate(mustTemplate(t, "second_template", "<div>Second</div>")).
		WithStyle(mustTemplate(t, "second_style", ".second { color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected second component, got error: %v", err)
	}

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyleList(
		[]Component{first, second},
		[]renderables.Params{
			{"color": "red"},
			{"color": "blue"},
		},
	)

	expected := ".second { color: blue; }"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildStyleWithMapParams(t *testing.T) {
	child, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("myComponent").
		WithTemplate(mustTemplate(t, "child_template", "<div>Child</div>")).
		WithStyle(mustTemplate(t, "child_style", ".myComponent { color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected child, got error: %v", err)
	}

	parent := mustComponent(t, "myParent", mustTemplate(t, "parent_template", "<section>{myComponent}</section>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyle(parent, renderables.Params{
		"myComponent": map[string]any{
			"color": "red",
		},
	})

	expected := ".myParent .myComponent { color: red; }"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildStyleWithRenderableParamsList(t *testing.T) {
	child, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("item").
		WithTemplate(mustTemplate(t, "item_template", "<li>Item</li>")).
		WithStyle(mustTemplate(t, "item_style", ".item { color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected child, got error: %v", err)
	}

	parent := mustComponent(t, "list", mustTemplate(t, "list_template", "<ul>{item}</ul>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyle(parent, renderables.Params{
		"item": []renderables.Params{
			{"color": "red"},
			{"color": "blue"},
		},
	})

	expected := ".list .item { color: red; }\n.list .item { color: blue; }"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildStyleWithMapParamsList(t *testing.T) {
	child, err := NewBuilder(renderables.NewStylableBuilder()).
		Create().
		WithKeyname("item").
		WithTemplate(mustTemplate(t, "item_template", "<li>Item</li>")).
		WithStyle(mustTemplate(t, "item_style", ".item { color: {color}; }")).
		Now()

	if err != nil {
		t.Fatalf("expected child, got error: %v", err)
	}

	parent := mustComponent(t, "list", mustTemplate(t, "list_template", "<ul>{item}</ul>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyle(parent, renderables.Params{
		"item": []map[string]any{
			{"color": "red"},
			{"color": "blue"},
		},
	})

	expected := ".list .item { color: red; }\n.list .item { color: blue; }"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestRendererRenderChildStyleWithDefaultValue(t *testing.T) {
	child := mustComponent(t, "value", mustTemplate(t, "value_template", "<span>Value</span>"))
	parent := mustComponent(t, "parent", mustTemplate(t, "parent_template", "<div>{value}</div>"), child)

	result := NewRenderer(templates.NewMustacheRenderer()).RenderStyle(parent, renderables.Params{
		"value": "raw-style",
	})

	expected := "raw-style"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
