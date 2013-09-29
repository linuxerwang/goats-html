package tests

import (
	"bytes"
	tmpl "goats-html/tests/templates/tag_dynamic_attrs_html"
	"testing"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagDynamicAttrs(t *testing.T) {
	args := &tmpl.TagDynamicAttrsTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagDynamicAttrsTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div class="test_tag">Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for rendering go:template tags, subsequent tag ==========

func TestTagDynamicAttrs1(t *testing.T) {
	args := &tmpl.TagDynamicAttrs1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagDynamicAttrs1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for call go:template tags ==========

func TestCallTagDynamicAttrs(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrs1(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrs1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrs1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrs2(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrs2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrs2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrs3(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrs3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrs3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrs4(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrs4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrs4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrs5(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrs5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrs5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrs6(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrs6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrs6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
