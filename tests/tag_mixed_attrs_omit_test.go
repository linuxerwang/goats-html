package tests

import (
	"bytes"
	tmpl "goats-html/tests/templates/tag_mixed_attrs_omit_html"
	"testing"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagMixedAttrsOmit(t *testing.T) {
	args := &tmpl.TagMixedAttrsOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagMixedAttrsOmitTemplate(&buffer, &settings)
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

func TestTagMixedAttrsOmit1(t *testing.T) {
	args := &tmpl.TagMixedAttrsOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagMixedAttrsOmit1Template(&buffer, &settings)
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

func TestCallTagMixedAttrsOmit(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsOmit1(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsOmit2(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsOmit2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsOmit2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsOmit3(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsOmit3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsOmit3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsOmit4(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsOmit4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsOmit4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsOmit5(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsOmit5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsOmit5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagMixedAttrsOmit6(t *testing.T) {
	args := &tmpl.CallTagMixedAttrsOmit6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagMixedAttrsOmit6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="test_tag tag1 caller_tag tag2">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
