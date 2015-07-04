package tests

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/tests/data"
	tmpl "github.com/linuxerwang/goats-html/tests/templates/replaces_html"
)

func TestWithoutReplace(t *testing.T) {
	blog := data.NewBlog()
	args := &tmpl.CallWithoutReplaceTemplateArgs{
		Blog: blog,
	}

	var buffer bytes.Buffer
	template := tmpl.NewCallWithoutReplaceTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div> Blog: <span>10010001</span><div> Slot1 suppressed. `+
			`</div> Hot Posts:<br><div> First hot post:<br> Id: <span>50001001</span><br> `+
			`Content: <span>I like eating donuts.</span><br> Second hot post:<br> Id: `+
			`<span>50001002</span><br> Content: <span>Jee, this morning I got up late.</span>`+
			`</div></div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}

func TestWithReplaces(t *testing.T) {
	blog := data.NewBlog()
	args := &tmpl.CallWithReplacesTemplateArgs{
		Blog: blog,
	}

	var buffer bytes.Buffer
	template := tmpl.NewCallWithReplacesTemplate(&buffer, &settings)
	err := template.Render(args)
	if err == nil {
		if buffer.String() != `<div><div> Blog: <span>10010001</span><div> Post Ids: `+
			`<span>50001001</span>, <span>50001002</span></div> Hot Posts:<br>`+
			`<div> Post1 content: <span>I like eating donuts.</span><br> `+
			`Post2 content: <span>Jee, this morning I got up late.</span></div></div></div>` {
			t.Error("Generated html: ", buffer.String())
		}
	} else {
		t.Error("Failed to render template. ", err)
	}
}
