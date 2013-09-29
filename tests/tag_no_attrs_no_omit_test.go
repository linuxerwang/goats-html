package tests

import (
	"bytes"
	tmpl "goats-html/tests/templates/tag_no_attrs_no_omit_html"
	"testing"
)

// ========== Tests for rendering go:template tags ==========

func TestTagNoAttrsNoOmit(t *testing.T) {
	args := &tmpl.TagNoAttrsNoOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagNoAttrsNoOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

// ========== Tests for rendering go:template tags, subsequent tag ==========

func TestTagNoAttrsNoOmit1(t *testing.T) {
	args := &tmpl.TagNoAttrsNoOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewTagNoAttrsNoOmit1Template(&buffer, &settings)
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

func TestCallTagNoAttrsNoOmit(t *testing.T) {
	args := &tmpl.CallTagNoAttrsNoOmitTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsNoOmitTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div>Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsNoOmit1(t *testing.T) {
	args := &tmpl.CallTagNoAttrsNoOmit1TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsNoOmit1Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsNoOmit2(t *testing.T) {
	args := &tmpl.CallTagNoAttrsNoOmit2TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsNoOmit2Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsNoOmit3(t *testing.T) {
	args := &tmpl.CallTagNoAttrsNoOmit3TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsNoOmit3Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsNoOmit4(t *testing.T) {
	args := &tmpl.CallTagNoAttrsNoOmit4TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsNoOmit4Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsNoOmit5(t *testing.T) {
	args := &tmpl.CallTagNoAttrsNoOmit5TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsNoOmit5Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div>Content</div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestCallTagNoAttrsNoOmit6(t *testing.T) {
	args := &tmpl.CallTagNoAttrsNoOmit6TemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewCallTagNoAttrsNoOmit6Template(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div class="caller_tag tag1">Content</div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
