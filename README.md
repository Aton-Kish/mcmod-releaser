# mcmod-releaser

## Installation

Download the pre-compiled binaries from the [releases page](https://github.com/Aton-Kish/mcmod-releaser/releases).

## Usage

```shell
export CURSEFORGE_API_TOKEN=<CurseForge API Token>

mcmod-releaser curseforge \
    --file ./path/to/example-mod-1.2.3.jar \
    --project-id <CurseForge Project ID> \
    --release-type release \
    --version 1.2.3 \
    --name "Example Mod v1.2.3" \
    --changelog "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)" \
    --environments "Client,Server" \
    --loaders "Fabric" \
    --java-versions "Jave 21" \
    --game-versions "1.21.9,1.21.10" \
    --dependencies "fabric-api=required,cloth-config=embedded,modmenu=optional"
```

## Changelog

Refer to the [CHANGELOG](./CHANGELOG.md).

## License

The syncup is licensed under the MIT License, see [LICENSE](./LICENSE).
