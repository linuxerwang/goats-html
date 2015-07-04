package tests

import (
	"bytes"
	"testing"

	tmpl "github.com/linuxerwang/goats-html/tests/templates/tag_static_attrs_omit_html"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagStaticAttrsOmit(t *testing.T) {
	args := &tmpl.TagStaticAttrsOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagStaticAttrsOmitTemplate(&buffer, &settings)
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

func TestTagStaticAttrsOmit1(t *testing.T) {
	args := &tmpl.TagStaticAttrsOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagStaticAttrsOmit1Template(&buffer, &settings)
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

func TestCallTagStaticAttrsOmit(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsOmit1(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsOmit2(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsOmit2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsOmit2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsOmit3(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsOmit3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsOmit3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsOmit4(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsOmit4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsOmit4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsOmit5(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsOmit5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsOmit5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsOmit6(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsOmit6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsOmit6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
