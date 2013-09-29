package runtime

import (
	"fmt"
	"io"
)

type LoopVar struct {
	Total       int
	Counter     int
	Counter0    int
	RevCounter  int
	RevCounter0 int
	First       bool
	Last        bool
}

type TemplateSettings struct {
	OmitDocType bool
}

type CallerAttrsFunc func() (TagAttrs, bool, bool)
type ReplaceableFunc func()

type Template interface {
	GetCallerAttrsFunc() CallerAttrsFunc
	SetCallerAttrsFunc(CallerAttrsFunc)
}

type BaseTemplate struct {
	writer          io.Writer
	settings        *TemplateSettings
	callerAttrsFunc CallerAttrsFunc
}

func (bt *BaseTemplate) WriteString(text string) {
	io.WriteString(bt.writer, text)
}

func (bt *BaseTemplate) FormatString(text string, a ...interface{}) {
	io.WriteString(bt.writer, fmt.Sprintf(text, a...))
}

func (bt *BaseTemplate) GetWriter() io.Writer {
	return bt.writer
}

func (bt *BaseTemplate) GetSettings() *TemplateSettings {
	return bt.settings
}

func (bt *BaseTemplate) GetCallerAttrsFunc() CallerAttrsFunc {
	return bt.callerAttrsFunc
}

func (bt *BaseTemplate) SetCallerAttrsFunc(callerAttrsFunc CallerAttrsFunc) {
	bt.callerAttrsFunc = callerAttrsFunc
}

func NewBaseTemplate(writer io.Writer, settings *TemplateSettings) *BaseTemplate {
	return &BaseTemplate{
		writer:   writer,
		settings: settings,
	}
}
