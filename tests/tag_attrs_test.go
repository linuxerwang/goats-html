package tests

import (
	"bytes"
	tmpl "goats-html/tests/templates/tag_attrs_html"
	"testing"
)


func TestTagAttrsInLoop(t *testing.T) {
	args := &tmpl.TagAttrsInLoopTemplateArgs{
		Names: []string{
			"John",
			"Matt",
			"Thomas",
			"Zoey",
		},
	}

	var buffer bytes.Buffer
	template := tmpl.NewTagAttrsInLoopTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		expected := `<div>` +
			`<input id="John" value="John" type="radio">` +
			`<label for="John">John</label><br>` +
			`<input id="Matt" value="Matt" type="radio">` +
			`<label for="Matt">Matt</label><br>` +
			`<input id="Thomas" value="Thomas" type="radio">` +
			`<label for="Thomas">Thomas</label><br>` +
			`<input id="Zoey" value="Zoey" type="radio">` +
			`<label for="Zoey">Zoey</label>` +
			`</div>`
		if buffer.String() != expected {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
