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
	"runtime/debug"
)

const (
	AppName = "mcmod-releaser"

	appVersionUnknown = "unknown"
)

var (
	version   = appVersionUnknown
	gitCommit = appVersionUnknown
	goVersion = appVersionUnknown
	goOS      = appVersionUnknown
	goArch    = appVersionUnknown
)

type AppVersion struct {
	Version   string `json:"version"`
	GitCommit string `json:"commit"`
	GoVersion string `json:"go"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

func NewAppVersion() *AppVersion {
	v := &AppVersion{
		Version:   version,
		GitCommit: gitCommit,
		GoVersion: goVersion,
		OS:        goOS,
		Arch:      goArch,
	}

	if v.Version == appVersionUnknown {
		v.loadFromBuildInfo()
	}

	return v
}

func (v *AppVersion) loadFromBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	v.Version = info.Main.Version
	v.GoVersion = info.GoVersion

	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			v.GitCommit = s.Value

			// Shorten the commit hash
			if len(v.GitCommit) > 7 {
				v.GitCommit = v.GitCommit[:7]
			}
		case "GOOS":
			v.OS = s.Value
		case "GOARCH":
			v.Arch = s.Value
		}
	}
}

type AppConfig struct {
	CurseForgeAPIToken string `env:"CURSEFORGE_API_TOKEN,required"`
}
