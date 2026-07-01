package templates

import (
	"testing"
)

func TestMustacheBuilderNow(t *testing.T) {
	template, err := NewMustacheBuilder().
		Create().
		WithKeyname("welcome").
		WithCode("Hello {name}, welcome to {project}.").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if template.Keyname() != "welcome" {
		t.Fatalf("expected keyname welcome, got %s", template.Keyname())
	}

	if template.Code() != "Hello {name}, welcome to {project}." {
		t.Fatalf("unexpected code: %s", template.Code())
	}

	if !template.HasParams() {
		t.Fatal("expected template to have params")
	}

	expectedParams := []string{"name", "project"}
	params := template.Params()

	if len(params) != len(expectedParams) {
		t.Fatalf("expected %d params, got %d", len(expectedParams), len(params))
	}

	for i, expected := range expectedParams {
		if params[i] != expected {
			t.Fatalf("expected param at index %d to be %s, got %s", i, expected, params[i])
		}
	}

	if !template.Contains("name") {
		t.Fatal("expected template to contain name")
	}

	if !template.Contains("project") {
		t.Fatal("expected template to contain project")
	}

	if template.Contains("missing") {
		t.Fatal("expected template to not contain missing")
	}
}

func TestMustacheBuilderTrimsKeyname(t *testing.T) {
	template, err := NewMustacheBuilder().
		Create().
		WithKeyname("  welcome  ").
		WithCode("Hello").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if template.Keyname() != "welcome" {
		t.Fatalf("expected trimmed keyname welcome, got %q", template.Keyname())
	}
}

func TestMustacheBuilderNowReturnsErrorWhenKeynameIsEmpty(t *testing.T) {
	_, err := NewMustacheBuilder().
		Create().
		WithKeyname("   ").
		WithCode("Hello {name}").
		Now()

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "template keyname is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustacheBuilderNowReturnsErrorWhenCodeIsEmpty(t *testing.T) {
	_, err := NewMustacheBuilder().
		Create().
		WithKeyname("welcome").
		WithCode("").
		Now()

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "template code is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustacheBuilderExtractsUniqueParams(t *testing.T) {
	template, err := NewMustacheBuilder().
		Create().
		WithKeyname("duplicate").
		WithCode("{name} {name} {email} {name}").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	params := template.Params()
	expected := []string{"name", "email"}

	if len(params) != len(expected) {
		t.Fatalf("expected %d params, got %d", len(expected), len(params))
	}

	for i := range expected {
		if params[i] != expected[i] {
			t.Fatalf("expected param %s at index %d, got %s", expected[i], i, params[i])
		}
	}
}

func TestMustacheBuilderIgnoresInvalidParams(t *testing.T) {
	template, err := NewMustacheBuilder().
		Create().
		WithKeyname("invalid").
		WithCode("{valid} {123invalid} {-invalid} {also_valid_123} {}").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	params := template.Params()
	expected := []string{"valid", "also_valid_123"}

	if len(params) != len(expected) {
		t.Fatalf("expected %d params, got %d: %#v", len(expected), len(params), params)
	}

	for i := range expected {
		if params[i] != expected[i] {
			t.Fatalf("expected param %s at index %d, got %s", expected[i], i, params[i])
		}
	}
}

func TestMustacheBuilderTemplateWithoutParams(t *testing.T) {
	template, err := NewMustacheBuilder().
		Create().
		WithKeyname("plain").
		WithCode("Hello world").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if template.HasParams() {
		t.Fatal("expected template to have no params")
	}

	if len(template.Params()) != 0 {
		t.Fatalf("expected no params, got %#v", template.Params())
	}
}

func TestNewMustacheRenderer(t *testing.T) {
	renderer := NewMustacheRenderer()

	if renderer == nil {
		t.Fatal("expected renderer")
	}
}

func TestMustacheRendererRender(t *testing.T) {
	template, err := NewMustacheBuilder().
		Create().
		WithKeyname("welcome").
		WithCode("Hello {name}, welcome to {project}.").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := NewMustacheRenderer().Render(template, map[string]string{
		"name":    "Steve",
		"project": "EventFlow",
	})

	expected := "Hello Steve, welcome to EventFlow."

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestMustacheRendererRenderMissingValueUsesEmptyString(t *testing.T) {
	template, err := NewMustacheBuilder().
		Create().
		WithKeyname("welcome").
		WithCode("Hello {name}, welcome to {project}.").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := NewMustacheRenderer().Render(template, map[string]string{
		"name": "Steve",
	})

	expected := "Hello Steve, welcome to ."

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestMustacheRendererRenderRepeatedParam(t *testing.T) {
	template, err := NewMustacheBuilder().
		Create().
		WithKeyname("repeat").
		WithCode("{name} likes {name}'s project").
		Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := NewMustacheRenderer().Render(template, map[string]string{
		"name": "Steve",
	})

	expected := "Steve likes Steve's project"

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
