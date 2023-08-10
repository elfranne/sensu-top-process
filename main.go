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
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-cpu-usage",
			Short:    "Check CPU usage and provide metrics",
			Keyspace: "sensu.io/plugins/check-cpu-usage/config",
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
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	if plugin.CPU == 100 {
		return sensu.CheckStateWarning, fmt.Errorf("that's just stupid")
	}
	if plugin.Memory == 100 {
		return sensu.CheckStateWarning, fmt.Errorf("that's just stupid")
	}
	if plugin.Scheme == "" {
		return sensu.CheckStateWarning, fmt.Errorf("scheme is required")
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
	process, _ := process.Processes()
	for _, p := range process {
		cpu, _ := p.CPUPercent()
		memory, _ := p.MemoryPercent()
		name, _ := p.Name()
		expanded := ExpandName(name, p)

		if cpu >= plugin.CPU || memory >= plugin.Memory {
			fmt.Printf("%s.process.cpu_percent.%s %f %d\n", plugin.Scheme, re.ReplaceAllString(expanded, "_"), Round(cpu, 0.1), time.Now().Unix())
			fmt.Printf("%s.process.memory_percent.%s %f %d\n", plugin.Scheme, re.ReplaceAllString(expanded, "_"), Round(float64(memory), 0.1), time.Now().Unix())
		}
	}
	return sensu.CheckStateOK, nil
}
