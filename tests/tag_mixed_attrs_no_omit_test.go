package tests

import (
	"bytes"
	"testing"

	tmpl "github.com/linuxerwang/goats-html/tests/templates/tag_mixed_attrs_no_omit_html"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagMixedAttrsNoOmit(t *testing.T) {
	args := &tmpl.TagMixedAttrsNoOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagMixedAttrsNoOmitTemplate(&buffer, &settings)
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

func TestTagMixedAttrsNoOmit1(t *testing.T) {
	args := &tmpl.TagMixedAttrsNoOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagMixedAttrsNoOmit1Template(&buffer, &settings)
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

func TestCallTagMixedAttrsNoOmit(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsNoOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsNoOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsNoOmit1(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsNoOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsNoOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsNoOmit2(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsNoOmit2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsNoOmit2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsNoOmit3(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsNoOmit3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsNoOmit3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsNoOmit4(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsNoOmit4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsNoOmit4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsNoOmit5(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsNoOmit5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsNoOmit5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsNoOmit6(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsNoOmit6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsNoOmit6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
