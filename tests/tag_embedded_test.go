package tests

import (
	"bytes"
	tmpl "goats-html/tests/templates/tag_embedded_html"
	"testing"
)

func TestTagEmbedded(t *testing.T) {
	expected := `<div class="class1"><div class="class2" style="padding: 5px;">` +
		`<span class="class3" id="a"> test </span></div></div>`
	args := &tmpl.EmbeddedTemplateArgs{}

	var buffer bytes.Buffer
	template := tmpl.NewEmbeddedTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != expected {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
