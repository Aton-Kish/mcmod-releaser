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

package registry

import (
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"

	"github.com/Aton-Kish/mcmod-releaser/internal/curseforge"
	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/command"
	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/model"
	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/repository"
	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/usecase"
)

type Registry struct {
	AppVersion *model.AppVersion
	AppConfig  *model.AppConfig

	CurseForgeRepository repository.CurseForgeRepository

	CurseForgeUseCase usecase.CurseForgeUseCase

	RootCommand       *cobra.Command
	VersionCommand    *cobra.Command
	CurseForgeCommand *cobra.Command
}

func New() (_ *Registry, err error) {
	defer model.WrapAppError(&err)

	cfg, err := env.ParseAs[model.AppConfig]()
	if err != nil {
		return nil, err
	}

	reg := &Registry{
		AppVersion: model.NewAppVersion(),
		AppConfig:  &cfg,
	}

	if err := reg.initRepository(); err != nil {
		return nil, err
	}

	if err := reg.initUseCase(); err != nil {
		return nil, err
	}

	if err := reg.initCommand(); err != nil {
		return nil, err
	}

	return reg, nil
}

func (r *Registry) initRepository() error {
	curseForgeClient, err := curseforge.NewClient(r.AppConfig.CurseForgeAPIToken)
	if err != nil {
		return err
	}

	r.CurseForgeRepository = repository.NewCurseForgeRepository(curseForgeClient)

	return nil
}

func (r *Registry) initUseCase() error {
	r.CurseForgeUseCase = usecase.NewCurseForgeUseCase(r.CurseForgeRepository)

	return nil
}

func (r *Registry) initCommand() error {
	r.RootCommand = command.NewRootCommand(r.AppVersion)
	r.VersionCommand = command.NewVersionCommand(r.AppVersion)
	r.CurseForgeCommand = command.NewCurseForgeCommand(r.CurseForgeUseCase)

	return nil
}
