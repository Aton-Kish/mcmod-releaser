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
	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/usecase"
)

type curseForgeCommandFlags struct {
	filePath string

	projectID    string
	releaseType  string
	version      string
	name         string
	changelog    string
	environments []string
	loaders      []string
	javaVersions []string
	gameVersions []string
	dependencies map[string]string

	parsedReleaseType  model.ModReleaseType
	parsedDependencies map[string]model.ModDependencyType
}

func (f *curseForgeCommandFlags) register(cmd *cobra.Command) error {
	cmd.Flags().StringVar(&f.filePath, "file", "", "The path to the mod file.")
	if err := cmd.MarkFlagFilename("file"); err != nil {
		return err
	}
	if err := cmd.MarkFlagRequired("file"); err != nil {
		return err
	}

	cmd.Flags().StringVar(&f.projectID, "project-id", "", "The CurseForge project ID.")
	if err := cmd.MarkFlagRequired("project-id"); err != nil {
		return err
	}

	cmd.Flags().StringVar(&f.releaseType, "release-type", "", "The release type. Allowed values: `release`, `beta`, `alpha`.")
	if err := cmd.MarkFlagRequired("release-type"); err != nil {
		return err
	}

	cmd.Flags().StringVar(&f.version, "version", "", "The mod version.")
	if err := cmd.MarkFlagRequired("version"); err != nil {
		return err
	}

	cmd.Flags().StringVar(&f.name, "name", "", "The mod name.")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		return err
	}

	cmd.Flags().StringVar(&f.changelog, "changelog", "", "The mod changelog. The changelog must be in Markdown format.")
	if err := cmd.MarkFlagRequired("changelog"); err != nil {
		return err
	}

	cmd.Flags().StringSliceVar(&f.environments, "environments", []string{}, "The mod environments.")

	cmd.Flags().StringSliceVar(&f.loaders, "loaders", []string{}, "The mod loaders.")
	if err := cmd.MarkFlagRequired("loaders"); err != nil {
		return err
	}

	cmd.Flags().StringSliceVar(&f.javaVersions, "java-versions", []string{}, "The mod Java versions.")

	cmd.Flags().StringSliceVar(&f.gameVersions, "game-versions", []string{}, "The mod game versions.")
	if err := cmd.MarkFlagRequired("game-versions"); err != nil {
		return err
	}

	cmd.Flags().StringToStringVar(&f.dependencies, "dependencies", map[string]string{}, "The mod dependencies. The key is the mod slug and the value is the dependency type. Allowed values for the dependency type: `required`, `embedded`, `optional`, `incompatible`, `tool`.")

	return nil
}

func (f *curseForgeCommandFlags) parse() error {
	if err := f.parseReleaseType(); err != nil {
		return err
	}

	if err := f.parseDependencies(); err != nil {
		return err
	}

	return nil
}

func (f *curseForgeCommandFlags) parseReleaseType() error {
	releaseType, err := model.NewModReleaseType(f.releaseType)
	if err != nil {
		return err
	}
	f.parsedReleaseType = releaseType
	return nil
}

func (f *curseForgeCommandFlags) parseDependencies() error {
	dependencies := make(map[string]model.ModDependencyType)
	for k, v := range f.dependencies {
		dependencyType, err := model.NewModDependencyType(v)
		if err != nil {
			return err
		}
		dependencies[k] = dependencyType
	}
	f.parsedDependencies = dependencies
	return nil
}

func NewCurseForgeCommand(curseForgeUseCase usecase.CurseForgeUseCase, optFns ...OptionFunc) *cobra.Command {
	opts := newOptions(optFns...)
	flags := new(curseForgeCommandFlags)
	cmd := &cobra.Command{
		Use:   "curseforge",
		Short: "Release a mod to CurseForge.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer model.WrapAppError(&err)

			if err := flags.parse(); err != nil {
				return err
			}

			output, err := curseForgeUseCase.Execute(cmd.Context(), &usecase.CurseForgeInput{
				FilePath: flags.filePath,
				Mod: &model.Mod{
					ProjectID:    flags.projectID,
					ReleaseType:  flags.parsedReleaseType,
					Version:      flags.version,
					Name:         flags.name,
					Changelog:    flags.changelog,
					Environments: flags.environments,
					Loaders:      flags.loaders,
					JavaVersions: flags.javaVersions,
					GameVersions: flags.gameVersions,
					Dependencies: flags.parsedDependencies,
				},
			})
			if err != nil {
				return err
			}

			data, err := json.Marshal(output)
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
	_ = flags.register(cmd)

	cmd.SetIn(opts.stdio.in)
	cmd.SetOut(opts.stdio.err)
	cmd.SetErr(opts.stdio.err)

	return cmd
}
