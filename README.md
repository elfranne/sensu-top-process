[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/elfranne/sensu-top-process)
![Go Test](https://github.com/elfranne/sensu-top-process/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/elfranne/sensu-top-process/workflows/goreleaser/badge.svg)

# sensu-top-process

## Table of Contents

- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Usage examples

### Basic Usage

To run the check with default parameters:

```shell
sensu-top-process
```

This will execute the check with the default CPU and memory thresholds (10%) and no specific scheme or expansion for process names.

### Custom Thresholds and Scheme

To specify custom thresholds and a scheme:

```shell
sensu-top-process --cpu 15.5 --memory 20 --scheme my_custom_scheme
```

This sets the CPU threshold to 15.5%, the memory threshold to 20%, and prepends `my_custom_scheme` to all emitted metrics.

### Expanding Process Names

To expand the process name to include arguments:

```shell
sensu-top-process --expand bash
```

This expands the names of processes named 'bash' to include their command-line arguments.

## Configuration

The `sensu-top-process` check can be configured with various command-line arguments:

- **`--cpu` or `-c`**: Set the CPU usage threshold as a percentage.
- **`--memory` or `-m`**: Set the memory usage threshold as a percentage.
- **`--scheme` or `-s`**: Specify a scheme to prepend to metric outputs.
- **`--expand` or `-e`**: Expand process name to include arguments (useful for processes like bash or powershell).

_Note: Detailed descriptions and default values for these configurations are provided in the [Overview](#overview) section._

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add elfranne/sensu-top-process
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/elfranne/sensu-top-process].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-top-process
  namespace: default
spec:
  command: sensu-top-process --cpu 15.5 --memory 20 --scheme my_scheme --expand bash
  subscriptions:
    - system
  runtime_assets:
    - elfranne/sensu-top-process
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-top-process repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/elfranne/sensu-top-process/blob/master/.github/workflows/release.yml
[5]: https://github.com/elfranne/sensu-top-process/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu/sensu-top-process/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
