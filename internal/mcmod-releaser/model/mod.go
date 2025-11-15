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

package model

import (
	"fmt"
)

type Mod struct {
	ID           string                       `json:"id"`
	ProjectID    string                       `json:"projectId"`
	ReleaseType  ModReleaseType               `json:"releaseType"`
	Version      string                       `json:"version"`
	Name         string                       `json:"name"`
	Changelog    string                       `json:"changelog"`
	Environments []string                     `json:"environments"`
	Loaders      []string                     `json:"loaders"`
	JavaVersions []string                     `json:"javaVersions"`
	GameVersions []string                     `json:"gameVersions"`
	Dependencies map[string]ModDependencyType `json:"dependencies"`
}

type ModReleaseType string

const (
	ModReleaseTypeRelease ModReleaseType = "release"
	ModReleaseTypeBeta    ModReleaseType = "beta"
	ModReleaseTypeAlpha   ModReleaseType = "alpha"
)

func NewModReleaseType(v string) (ModReleaseType, error) {
	switch v {
	case "release":
		return ModReleaseTypeRelease, nil
	case "beta":
		return ModReleaseTypeBeta, nil
	case "alpha":
		return ModReleaseTypeAlpha, nil
	default:
		return "", fmt.Errorf("invalid release type: `%s`. allowed values: `release`, `beta`, `alpha`", v)
	}
}

type ModDependencyType string

const (
	ModDependencyTypeRequired     ModDependencyType = "required"
	ModDependencyTypeOptional     ModDependencyType = "optional"
	ModDependencyTypeIncompatible ModDependencyType = "incompatible"
	ModDependencyTypeEmbedded     ModDependencyType = "embedded"
	ModDependencyTypeTool         ModDependencyType = "tool"
)

func NewModDependencyType(v string) (ModDependencyType, error) {
	switch v {
	case "required":
		return ModDependencyTypeRequired, nil
	case "optional":
		return ModDependencyTypeOptional, nil
	case "incompatible":
		return ModDependencyTypeIncompatible, nil
	case "embedded":
		return ModDependencyTypeEmbedded, nil
	case "tool":
		return ModDependencyTypeTool, nil
	default:
		return "", fmt.Errorf("invalid dependency type: `%s`. allowed values: `required`, `optional`, `incompatible`, `embedded`, `tool`", v)
	}
}
