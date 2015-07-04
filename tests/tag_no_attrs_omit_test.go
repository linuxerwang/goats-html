package tests

import (
	"bytes"
	"testing"

	tmpl "github.com/linuxerwang/goats-html/tests/templates/tag_no_attrs_omit_html"
)

// ========== Tests for rendering go:template tags, first tag ==========

func TestTagNoAttrsOmit(t *testing.T) {
	args := &tmpl.TagNoAttrsOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagNoAttrsOmitTemplate(&buffer, &settings)
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

func TestTagNoAttrsOmit1(t *testing.T) {
	args := &tmpl.TagNoAttrsOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagNoAttrsOmit1Template(&buffer, &settings)
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

func TestCallTagNoAttrsOmit(t *testing.T) {
	args := &tmpl.CallTagNoAttrsOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsOmit1(t *testing.T) {
	args := &tmpl.CallTagNoAttrsOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsOmit2(t *testing.T) {
	args := &tmpl.CallTagNoAttrsOmit2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsOmit2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsOmit3(t *testing.T) {
	args := &tmpl.CallTagNoAttrsOmit3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsOmit3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsOmit4(t *testing.T) {
	args := &tmpl.CallTagNoAttrsOmit4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsOmit4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsOmit5(t *testing.T) {
	args := &tmpl.CallTagNoAttrsOmit5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsOmit5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsOmit6(t *testing.T) {
	args := &tmpl.CallTagNoAttrsOmit6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsOmit6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
