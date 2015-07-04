package tests

import (
	"bytes"
	"testing"

	tmpl "github.com/linuxerwang/goats-html/tests/templates/tag_dynamic_attrs_omit_html"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagDynamicAttrsOmit(t *testing.T) {
	args := &tmpl.TagDynamicAttrsOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagDynamicAttrsOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `Content` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for rendering go:template tags, subsequent tag ==========

func TestTagDynamicAttrsOmit1(t *testing.T) {
	args := &tmpl.TagDynamicAttrsOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagDynamicAttrsOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for call go:template tags ==========

func TestCallTagDynamicAttrsOmit(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsOmit1(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsOmit2(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsOmit2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsOmit2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsOmit3(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsOmit3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsOmit3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsOmit4(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsOmit4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsOmit4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsOmit5(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsOmit5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsOmit5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsOmit6(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsOmit6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsOmit6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
