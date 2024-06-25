package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParamsGetEnvVar struct {
	path string
}

type ExpectedGetEnvVar struct {
	str   string
	idDel bool
	err   error
}

func TestGetEnvVar(t *testing.T) {
	cases := []struct {
		name     string
		params   ParamsGetEnvVar
		expected ExpectedGetEnvVar
	}{
		{
			name: "BAR_Ok",
			params: ParamsGetEnvVar{
				path: "./testdata/env/BAR",
			},
			expected: ExpectedGetEnvVar{
				str:   "bar",
				idDel: false,
				err:   error(nil),
			},
		},
		{
			name: "EMPTY_Ok",
			params: ParamsGetEnvVar{
				path: "./testdata/env/EMPTY",
			},
			expected: ExpectedGetEnvVar{
				str:   "",
				idDel: true,
				err:   error(nil),
			},
		},
		{
			name: "FOO_Ok",
			params: ParamsGetEnvVar{
				path: "./testdata/env/FOO",
			},
			expected: ExpectedGetEnvVar{
				str:   "   foo\nwith new line",
				idDel: false,
				err:   error(nil),
			},
		},
		{
			name: "HELLO_Ok",
			params: ParamsGetEnvVar{
				path: "./testdata/env/HELLO",
			},
			expected: ExpectedGetEnvVar{
				str:   `"hello"`,
				idDel: false,
				err:   error(nil),
			},
		},
		{
			name: "UNSET_Ok",
			params: ParamsGetEnvVar{
				path: "./testdata/env/UNSET",
			},
			expected: ExpectedGetEnvVar{
				str:   "",
				idDel: true,
				err:   error(nil),
			},
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			res, isDelRes, errRes := GetEnvVar(cs.params.path)
			assert.ErrorIs(t, errRes, cs.expected.err)
			assert.Equal(t, res, cs.expected.str)
			assert.Equal(t, isDelRes, cs.expected.idDel)
		})
	}
}

type ParamsReadDir struct {
	path string
}

type ExpectedReadDir struct {
	env Environment
	err error
}

func TestReadDir(t *testing.T) {
	cases := []struct {
		name     string
		params   ParamsReadDir
		expected ExpectedReadDir
	}{
		{
			name: "_Ok",
			params: ParamsReadDir{
				path: "./testdata/env/",
			},
			expected: ExpectedReadDir{
				env: Environment{
					"BAR": EnvValue{
						Value:      "bar",
						NeedRemove: false,
					},
					"EMPTY": EnvValue{
						Value:      "",
						NeedRemove: true,
					},
					"FOO": EnvValue{
						Value:      "   foo\nwith new line",
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
				err: error(nil),
			},
		},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			envRes, errRes := ReadDir(cs.params.path)
			assert.ErrorIs(t, errRes, cs.expected.err)
			assert.Equal(t, envRes, cs.expected.env)
		})
	}
}
