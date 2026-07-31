[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/elfranne/sensu-top-process)
![Go Test](https://github.com/elfranne/sensu-top-process/workflows/Go%20Test/badge.svg)
![Go Lint](https://github.com/elfranne/sensu-top-process/workflows/Go%20Lint/badge.svg)
![goreleaser](https://github.com/elfranne/sensu-top-process/workflows/goreleaser/badge.svg)

# sensu-top-process

## Table of Contents

- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

`sensu-top-process` is a [Sensu Check][1] that emits Graphite plaintext metrics for the heaviest
processes on a host. Every process whose CPU usage **or** memory usage is at or above the
configured threshold produces two metrics:

```text
<scheme>.process.cpu_percent.<process_name> <value> <timestamp>
<scheme>.process.memory_percent.<process_name> <value> <timestamp>
```

Values are percentages rounded to one decimal place, and the timestamp is Unix epoch seconds.
Characters that are not valid in a Graphite metric path (`-`, whitespace, `/`, `:`, `.`, `,`, `=`)
become `_` in the process name. A *run* of the same character collapses into a single underscore
— `my--app` becomes `my_app` and `python3.11.2` becomes `python3_11_2` — but adjacent characters
from different classes each contribute one, so `/usr/bin/python3 --verbose` becomes
`_usr_bin_python3__verbose`.

Example output:

```text
my_scheme.process.cpu_percent.firefox 23.400000 1738310400
my_scheme.process.memory_percent.firefox 12.100000 1738310400
```

The check always exits `0` (OK) once its arguments validate — it is a metrics collector, not an
alerting check. Configure alerting on the resulting metrics instead.

Process data is collected through [gopsutil][2], so Linux, macOS and Windows are all supported.
Prebuilt assets are published for linux (amd64, 386, arm64, armv6, armv7), darwin/amd64 and
windows/amd64.

## Usage examples

### Basic usage

`--scheme` is required; the check exits with a warning if it is not set.

```shell
sensu-top-process --scheme my_scheme
```

This reports every process at or above the default thresholds of 10% CPU or 10% memory, prefixing
each metric with `my_scheme`. The run takes about a second, because CPU usage is measured over a
sample window rather than read from a single instant.

### Custom thresholds

```shell
sensu-top-process --scheme my_scheme --cpu 15.5 --memory 20
```

This sets the CPU threshold to 15.5% and the memory threshold to 20%. A process is reported if it
crosses *either* threshold.

### Changing the sample window

```shell
sensu-top-process --scheme my_scheme --sample 5
```

CPU usage is measured by reading each process's CPU time twice, `--sample` seconds apart, and
reporting the difference. A longer window smooths out brief spikes and makes the numbers steadier;
a shorter one returns faster but reads more coarsely. The check sleeps once per run, not once per
process, so the added time is the sample duration no matter how many processes are running.

### Expanding process names

Interpreter processes all share the same name (`bash`, `python`, `powershell`…), which makes their
metrics collide. Use `--expand` to replace the name of a single process with its full command line:

```shell
sensu-top-process --scheme my_scheme --expand bash
```

Processes named `bash` are then reported as, for example,
`my_scheme.process.cpu_percent._usr_bin_bash__home_user_backup_sh`. Only one process name can be
expanded per check.

### Help

```shell
sensu-top-process --help
```

## Configuration

| Argument | Short | Type | Default | Description |
| --- | --- | --- | --- | --- |
| `--cpu` | `-c` | float | `10` | Report processes at or above this CPU percentage. Must be above `0` and not exactly `100`. |
| `--memory` | `-m` | float | `10` | Report processes at or above this memory percentage. Must be above `0` and not exactly `100`. |
| `--scheme` | `-s` | string | *(none)* | Prefix prepended to every metric. **Required.** |
| `--expand` | `-e` | string | *(none)* | Process name to expand to its full command line. |
| `--sample` | | float | `1` | Seconds to measure CPU usage over. The check sleeps this long. Must be above `0`. |

Every argument can also be set through the check annotation keyspace
`sensu.io/plugins/check-cpu-usage/config` (for example the annotation
`sensu.io/plugins/check-cpu-usage/config/cpu: "25"`).

### Asset registration

[Sensu Assets][3] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```shell
sensuctl asset add elfranne/sensu-top-process
```

If you're using an earlier version of sensuctl, you can find the asset on the
[Bonsai Asset Index](https://bonsai.sensu.io/assets/elfranne/sensu-top-process).

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-top-process
  namespace: default
spec:
  command: sensu-top-process --cpu 15.5 --memory 20 --scheme my_scheme --expand bash --sample 1
  subscriptions:
    - system
  runtime_assets:
    - elfranne/sensu-top-process
  interval: 60
  timeout: 10
  publish: true
  output_metric_format: graphite_plaintext
  output_metric_handlers:
    - influxdb
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

Building requires the Go version declared in [`go.mod`](go.mod) (currently Go 1.26) or later.

From the local path of the sensu-top-process repository:

```shell
go build
go test ./...
```

## Additional notes

- **The check sleeps for `--sample` seconds on every run.** CPU usage is a rate, so it cannot be
  read from a single instant — the check takes one reading, waits, takes a second, and reports the
  difference. It sleeps once for the whole run rather than once per process, so the cost is the
  sample duration regardless of how many processes are on the host. If your check definition sets
  a `timeout`, make sure it is comfortably larger than `--sample` or the check will be killed
  before it reports.
- **Very short samples read coarsely.** Linux accounts process CPU time in clock ticks, typically
  10ms, so a 100ms sample counts only about ten of them and the results step between coarse
  values. A sample of `1` second or more gives stable numbers.
- On a multi-core host a process using more than one core reports above 100% (four saturated cores
  is roughly `400`), so CPU thresholds above 100 are meaningful. Memory percentages are of total
  system memory and stay within 0–100.
- Processes are enumerated once, before the sample. A process that starts during the sample window
  is not reported until the next run, and one that exits during it is dropped.
- Reading other users' processes may require elevated privileges; processes the agent cannot inspect
  are silently skipped.

## Contributing

For more information about contributing to this plugin, see [Contributing][4].

[1]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[2]: https://github.com/shirou/gopsutil
[3]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[4]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
