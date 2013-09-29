package tests

import (
	"goats-html/goats/runtime"
)

var settings runtime.TemplateSettings

func init() {
	runtime.InitGoats(nil)
	settings = runtime.TemplateSettings{
		OmitDocType: false,
	}
}
