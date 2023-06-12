package helpers

import (
	"bytes"
	"strings"
	"tddapps.com/truecrypt/internal"
	"tddapps.com/truecrypt/internal/paths"
)

type FakeInput struct {
	in  *internal.Input
	out bytes.Buffer
}

func NewFakeInput(lines ...string) *FakeInput {
	return NewFakeInputWithSettingsPath("should_not_be_used", lines...)
}

func NewFakeInputWithSettingsPath(settingsPath paths.FilePath, lines ...string) *FakeInput {
	f := FakeInput{}

	f.in = &internal.Input{
		IO: internal.IO{
			// always add a newline at the end
			// to properly mock one enter per line
			Reader: strings.NewReader(
				strings.Join(lines, "\n") + "\n",
			),
			Writer: &f.out,
		},
		SettingsPath: settingsPath,
	}

	return &f
}

func (f *FakeInput) In() *internal.Input {
	return f.in
}

func (f *FakeInput) Out() string {
	return f.out.String()
}
