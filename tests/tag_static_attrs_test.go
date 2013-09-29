package tests

import (
	"bytes"
	tmpl "goats-html/tests/templates/tag_static_attrs_html"
	"testing"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagStaticAttrs(t *testing.T) {
	args := &tmpl.TagStaticAttrsTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagStaticAttrsTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div class="tag1">Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for rendering go:template tags, subsequent tag ==========

func TestTagStaticAttrs1(t *testing.T) {
	args := &tmpl.TagStaticAttrs1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagStaticAttrs1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for call go:template tags ==========

func TestCallTagStaticAttrs(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrs1(t *testing.T) {
	args := &tmpl.CallTagStaticAttrs1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrs1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrs2(t *testing.T) {
	args := &tmpl.CallTagStaticAttrs2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrs2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrs3(t *testing.T) {
	args := &tmpl.CallTagStaticAttrs3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrs3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrs4(t *testing.T) {
	args := &tmpl.CallTagStaticAttrs4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrs4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrs5(t *testing.T) {
	args := &tmpl.CallTagStaticAttrs5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrs5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrs6(t *testing.T) {
	args := &tmpl.CallTagStaticAttrs6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrs6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
