package tests

import (
	"bytes"
	"testing"

	tmpl "github.com/linuxerwang/goats-html/tests/templates/tag_static_attrs_no_omit_html"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagStaticAttrsNoOmit(t *testing.T) {
	args := &tmpl.TagStaticAttrsNoOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagStaticAttrsNoOmitTemplate(&buffer, &settings)
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

func TestTagStaticAttrsNoOmit1(t *testing.T) {
	args := &tmpl.TagStaticAttrsNoOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagStaticAttrsNoOmit1Template(&buffer, &settings)
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

func TestCallTagStaticAttrsNoOmit(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsNoOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsNoOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsNoOmit1(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsNoOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsNoOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsNoOmit2(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsNoOmit2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsNoOmit2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsNoOmit3(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsNoOmit3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsNoOmit3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsNoOmit4(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsNoOmit4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsNoOmit4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsNoOmit5(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsNoOmit5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsNoOmit5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagStaticAttrsNoOmit6(t *testing.T) {
	args := &tmpl.CallTagStaticAttrsNoOmit6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagStaticAttrsNoOmit6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
