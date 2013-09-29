package tests

import (
	"bytes"
	tmpl "goats-html/tests/templates/tag_no_attrs_html"
	"testing"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagNoAttrs(t *testing.T) {
	args := &tmpl.TagNoAttrsTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagNoAttrsTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for call go:template tags, subsequent tag ==========

func TestTagNoAttrs1(t *testing.T) {
	args := &tmpl.TagNoAttrs1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagNoAttrs1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div>Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for call go:template tags ==========

func TestCallTagNoAttrs(t *testing.T) {
	args := &tmpl.CallTagNoAttrsTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div>Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrs1(t *testing.T) {
	args := &tmpl.CallTagNoAttrs1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrs1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrs2(t *testing.T) {
	args := &tmpl.CallTagNoAttrs2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrs2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrs3(t *testing.T) {
	args := &tmpl.CallTagNoAttrs3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrs3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrs4(t *testing.T) {
	args := &tmpl.CallTagNoAttrs4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrs4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrs5(t *testing.T) {
	args := &tmpl.CallTagNoAttrs5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrs5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrs6(t *testing.T) {
	args := &tmpl.CallTagNoAttrs6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrs6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
