package cmd

import (
	"fmt"
	"strings"
)

type NukeParameters struct {
	ConfigPath string
	Format     string
	OutputFile string

	Targets  []string
	Excludes []string

	NoDryRun   bool
	Force      bool
	ForceSleep int

	MaxWaitRetries int
}

func (p *NukeParameters) Validate() error {
	if strings.TrimSpace(p.ConfigPath) == "" {
		return fmt.Errorf("You have to specify the --config flag.\n")
	}

	p.Format = strings.TrimSpace(p.Format)
	switch p.Format {
	case "json":
		fallthrough
	case "yaml":
		fallthrough
	case "text":
		return nil
	default:
		return fmt.Errorf("Unknown output format\n")
	}
}
