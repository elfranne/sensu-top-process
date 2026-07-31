package main

import (
	"fmt"
	"math"
	"regexp"
	"time"

	corev2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/shirou/gopsutil/v3/process"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	CPU    float64
	Memory float32
	Scheme string
	Expand string
	Sample float64
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-top-process",
			Short:    "Check CPU usage and provide metrics",
			Keyspace: "sensu.io/plugins/sensu-top-process/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[float64]{
			Path:      "cpu",
			Argument:  "cpu",
			Shorthand: "c",
			Default:   float64(10),
			Usage:     "Show metrics for processes above CPU x%",
			Value:     &plugin.CPU,
		},
		&sensu.PluginConfigOption[float32]{
			Path:      "memory",
			Argument:  "memory",
			Shorthand: "m",
			Default:   float32(10),
			Usage:     "Show metrics for processes above Memory x%",
			Value:     &plugin.Memory,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "scheme",
			Argument:  "scheme",
			Shorthand: "s",
			Default:   "",
			Usage:     "Scheme to prepend metric",
			Value:     &plugin.Scheme,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "expand",
			Argument:  "expand",
			Shorthand: "e",
			Default:   "",
			Usage:     "Expand name for process to include argurment(s) (usefull for bash or powershell)",
			Value:     &plugin.Expand,
		},
		&sensu.PluginConfigOption[float64]{
			Path:     "sample",
			Argument: "sample",
			Default:  float64(1),
			Usage:    "Seconds to sample CPU usage over, the check sleeps this long",
			Value:    &plugin.Sample,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	if plugin.CPU <= 0 || plugin.CPU == 100 {
		return sensu.CheckStateWarning, fmt.Errorf("cpu %v is just stupid, use a value above 0 that is not 100", plugin.CPU)
	}
	if plugin.Memory <= 0 || plugin.Memory == 100 {
		return sensu.CheckStateWarning, fmt.Errorf("memory %v is just stupid, use a value above 0 that is not 100", plugin.Memory)
	}
	if plugin.Scheme == "" {
		return sensu.CheckStateWarning, fmt.Errorf("scheme is required")
	}
	if plugin.Sample <= 0 {
		return sensu.CheckStateWarning, fmt.Errorf("sample %v is just stupid, use a value above 0", plugin.Sample)
	}

	return sensu.CheckStateOK, nil
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func ExpandName(name string, p *process.Process) string {
	if plugin.Expand == name {
		cmd, _ := p.Cmdline()
		return cmd
	} else {
		return name
	}
}

func executeCheck(event *corev2.Event) (int, error) {
	re := regexp.MustCompile(`-+|\s+|/+|:+|\.+|,+|=+`)
	procs, _ := process.Processes()

	// Percent(0) reports CPU usage since the previous call on the same process,
	// so seed every process, sleep once, then read the delta back. Sleeping here
	// rather than passing an interval to Percent keeps the cost at one sample for
	// the whole run instead of one sample per process.
	for _, p := range procs {
		_, _ = p.Percent(0)
	}
	time.Sleep(time.Duration(plugin.Sample * float64(time.Second)))

	now := time.Now().Unix()
	for _, p := range procs {
		cpu, err := p.Percent(0)
		if err != nil {
			// Process exited during the sample.
			continue
		}
		memory, _ := p.MemoryPercent()
		name, _ := p.Name()
		expanded := ExpandName(name, p)

		if cpu >= plugin.CPU || memory >= plugin.Memory {
			fmt.Printf("%s.process.cpu_percent.%s %f %d\n", plugin.Scheme, re.ReplaceAllString(expanded, "_"), Round(cpu, 0.1), now)
			fmt.Printf("%s.process.memory_percent.%s %f %d\n", plugin.Scheme, re.ReplaceAllString(expanded, "_"), Round(float64(memory), 0.1), now)
		}
	}
	return sensu.CheckStateOK, nil
}
