package tests

import (
	"bytes"
	"testing"

	tmpl "github.com/linuxerwang/goats-html/tests/templates/tag_dynamic_attrs_no_omit_html"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagDynamicAttrsNoOmit(t *testing.T) {
	args := &tmpl.TagDynamicAttrsNoOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagDynamicAttrsNoOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div class="test_tag tag1">Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for rendering go:template tags, subsequent tag ==========

func TestTagDynamicAttrsNoOmit1(t *testing.T) {
	args := &tmpl.TagDynamicAttrsNoOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagDynamicAttrsNoOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for call go:template tags ==========

func TestCallTagDynamicAttrsNoOmit(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsNoOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsNoOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsNoOmit1(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsNoOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsNoOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsNoOmit2(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsNoOmit2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsNoOmit2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsNoOmit3(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsNoOmit3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsNoOmit3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsNoOmit4(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsNoOmit4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsNoOmit4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsNoOmit5(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsNoOmit5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsNoOmit5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagDynamicAttrsNoOmit6(t *testing.T) {
	args := &tmpl.CallTagDynamicAttrsNoOmit6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagDynamicAttrsNoOmit6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
