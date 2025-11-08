// Copyright (c) 2025 Aton-Kish
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package command

import (
	"errors"
	"io"
	"os"

	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/model"
)

type stdio struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

var (
	defaultStdio = stdio{in: os.Stdin, out: os.Stdout, err: os.Stderr}
)

type options struct {
	stdio stdio
}

type OptionFunc func(o *options)

func newOptions(optFns ...OptionFunc) *options {
	o := &options{
		stdio: defaultStdio,
	}

	for _, fn := range optFns {
		fn(o)
	}

	return o
}

func WithStdio(in io.Reader, out, err io.Writer) OptionFunc {
	return func(o *options) {
		o.stdio = stdio{in: in, out: out, err: err}
	}
}

func handleError(errp *error) {
	if errp == nil || *errp == nil {
		return
	}

	if appErr := new(model.AppError); !errors.As(*errp, &appErr) {
		*errp = model.NewAppError(model.AppErrorCodeUnexpectedError, "Unexpected error occurred.", *errp)
	}
}
