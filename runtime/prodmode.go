// +build !goats_devmod

package runtime

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

var (
	goatsSettings *GoatsSettings
)

func NewGoatsSettings() *GoatsSettings {
	return &GoatsSettings{}
}

func InitGoats(settings *GoatsSettings) {
	if settings == nil {
		goatsSettings = NewGoatsSettings()
	} else {
		goatsSettings = settings
	}
}

func DecodeRpcRequestOrFail(input io.Reader, settings *TemplateSettings, args interface{}) {
	decoder := gob.NewDecoder(input)
	err := decoder.Decode(settings)
	if err == nil {
		err = decoder.Decode(args)
	}
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(2)
	}
}
