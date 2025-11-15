# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<a name="unreleased"></a>
## [Unreleased]

<a name="v0.1.0"></a>
## v0.1.0 - 2025-11-15

### Bug fixes

- [`2b6f4b4`](https://github.com/Aton-Kish/mcmod-releaser/commit/2b6f4b43c24341efe9bed6b607d2d8a9b212aafc) added GOVERSION setting step
- [`333750e`](https://github.com/Aton-Kish/mcmod-releaser/commit/333750ead1e0a61c1ab4447d982d45d50713e763) removed required flag

### Chores

- [`151af18`](https://github.com/Aton-Kish/mcmod-releaser/commit/151af181fc2709ee1371875dbe16e1e2329a39a5) added release task
- [`c89058f`](https://github.com/Aton-Kish/mcmod-releaser/commit/c89058f56d86ffc5adf7960c82c36e31048e6e38) added test task
- [`f545e69`](https://github.com/Aton-Kish/mcmod-releaser/commit/f545e693143ff313dc8f41102a995f965c5bee64) installed govulncheck
- [`a80fb88`](https://github.com/Aton-Kish/mcmod-releaser/commit/a80fb8828bcf58ed9a094d93478c83819f00875b) initial commit

### Code refactoring

- [`3c5c927`](https://github.com/Aton-Kish/mcmod-releaser/commit/3c5c9277d9db4f7eeecb66ca369409436d550a01) `file (*os.File)` arg was replaced with `path (string)`
- [`9bf7836`](https://github.com/Aton-Kish/mcmod-releaser/commit/9bf7836ef67be15e8189d0799ead06a33ba05929) splitted CreateModCurseForgey.Repositor method

### Continuous integration

- [`8d9055b`](https://github.com/Aton-Kish/mcmod-releaser/commit/8d9055b85b420d4e69fb6bd8fcc66e1f9d42f0fb) added GitHub Actions workflows

### Documentation

- [`a80170b`](https://github.com/Aton-Kish/mcmod-releaser/commit/a80170bd0ee1fa4c0560aa8531bc5a90148cb697) updated README.md

### Features

- [`aa0d051`](https://github.com/Aton-Kish/mcmod-releaser/commit/aa0d0518c69554286b96cbce1e8a881beb9e8466) added curseforge command
- [`72da716`](https://github.com/Aton-Kish/mcmod-releaser/commit/72da716d5b5a087f768c7e6e9ba1322a3f9005a9) added curseforge repository
- [`bbb9c5d`](https://github.com/Aton-Kish/mcmod-releaser/commit/bbb9c5d576e547970dde1ce45be59fabb8ccf9b1) added curseforge api client
- [`8679768`](https://github.com/Aton-Kish/mcmod-releaser/commit/86797682d54544a02aef04ee53f4aae43d04fbf4) added initial cli impl

[Unreleased]: https://github.com/Aton-Kish/mcmod-releaser/compare/v0.1.0...HEAD
