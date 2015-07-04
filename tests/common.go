package tests

import (
	"github.com/linuxerwang/goats-html/runtime"
)

var settings runtime.TemplateSettings

func init() {
	runtime.InitGoats(nil)
	settings = runtime.TemplateSettings{
		OmitDocType: false,
	}
}
