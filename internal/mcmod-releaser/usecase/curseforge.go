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

package usecase

import (
	"context"

	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/model"
	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/repository"
)

type CurseForgeInput struct {
	FilePath string
	Mod      *model.Mod
}

type CurseForgeOutput struct {
	FilePath string     `json:"file"`
	Mod      *model.Mod `json:"mod"`
}

type CurseForgeUseCase interface {
	Execute(ctx context.Context, input *CurseForgeInput) (*CurseForgeOutput, error)
}

type curseForgeModCreateUseCase struct {
	curseForgeRepository repository.CurseForgeRepository
}

func NewCurseForgeUseCase(curseForgeRepository repository.CurseForgeRepository) CurseForgeUseCase {
	return &curseForgeModCreateUseCase{
		curseForgeRepository: curseForgeRepository,
	}
}

func (u *curseForgeModCreateUseCase) Execute(ctx context.Context, input *CurseForgeInput) (*CurseForgeOutput, error) {
	mod, err := u.curseForgeRepository.CreateMod(ctx, input.FilePath, input.Mod)
	if err != nil {
		return nil, err
	}

	return &CurseForgeOutput{
		FilePath: input.FilePath,
		Mod:      mod,
	}, nil
}
