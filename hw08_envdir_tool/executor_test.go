package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParamsRunCmd struct {
	cmd []string
	env Environment
}

type ExpectedRunCmd struct {
	code int
}

func TestRunCmd(t *testing.T) {
	cases := []struct {
		name     string
		params   ParamsRunCmd
		expected ExpectedRunCmd
	}{
		{
			name: "_Ok",
			params: ParamsRunCmd{
				cmd: []string{"ls", "-l", "-a"},
				env: Environment{
					"BAR": EnvValue{
						Value:      "bar",
						NeedRemove: false,
					},
					"EMPTY": EnvValue{
						Value:      " ",
						NeedRemove: false,
					},
					"FOO": EnvValue{
						Value:      "   foo",
						NeedRemove: false,
					},
					"HELLO": EnvValue{
						Value:      `"hello"`,
						NeedRemove: false,
					},
					"UNSET": EnvValue{
						Value:      "",
						NeedRemove: true,
					},
				},
			},
			expected: ExpectedRunCmd{
				code: 0,
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			codeRes := RunCmd(cs.params.cmd, cs.params.env)
			assert.Equal(t, codeRes, cs.expected.code)
		})
	}
}
