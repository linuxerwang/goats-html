package tests

import (
	"bytes"
	tmpl "goats-html/tests/templates/tag_mixed_attrs_html"
	"testing"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagMixedAttrs1(t *testing.T) {
	args := &tmpl.TagMixedAttrs1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagMixedAttrs1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for rendering go:template tags, subsequent tag ==========

func TestTagMixedAttrs(t *testing.T) {
	args := &tmpl.TagMixedAttrsTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagMixedAttrsTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div class="test_tag tag1">Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for call go:template tags ==========

func TestCallTagMixedAttrs(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrs1(t *testing.T) {
	args := &tmpl.CallTagMixedAttrs1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrs1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrs2(t *testing.T) {
	args := &tmpl.CallTagMixedAttrs2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrs2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrs3(t *testing.T) {
	args := &tmpl.CallTagMixedAttrs3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrs3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrs4(t *testing.T) {
	args := &tmpl.CallTagMixedAttrs4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrs4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrs5(t *testing.T) {
	args := &tmpl.CallTagMixedAttrs5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrs5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrs6(t *testing.T) {
	args := &tmpl.CallTagMixedAttrs6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrs6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
