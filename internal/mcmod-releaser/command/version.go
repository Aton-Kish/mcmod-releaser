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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/model"
)

func NewVersionCommand(version *model.AppVersion, optFns ...OptionFunc) *cobra.Command {
	opts := newOptions(optFns...)
	cmd := &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Display the %s version.", model.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer handleError(&err)

			data, err := json.Marshal(version)
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintln(opts.stdio.out, string(data)); err != nil {
				return err
			}

			return nil
		},
		SilenceUsage: true,
	}

	cmd.SetIn(opts.stdio.in)
	cmd.SetOut(opts.stdio.err)
	cmd.SetErr(opts.stdio.err)

	return cmd
}
