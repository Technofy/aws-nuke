package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/rebuy-de/aws-nuke/pkg/types"
	"gopkg.in/yaml.v2"
	"io"
	"sync"

	"github.com/rebuy-de/aws-nuke/resources"
)

type exportableItem struct {
	Region string `json:"region" yaml:"region"`
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
	State string `json:"state" yaml:"state"`
	Properties types.Properties `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type docExport struct {
	Items []exportableItem `json:"items" yaml:"items"`
}

type docExporterType int

const (
	docExporterJson docExporterType = iota
	docExporterYaml
)

type DocItemPrinter struct {
	exporterType docExporterType
	w io.Writer
	export docExport
	lock sync.Mutex
}

func NewJsonItemPrinter(w io.Writer) *DocItemPrinter {
	return &DocItemPrinter{
		w: w,
		exporterType: docExporterJson,
	}
}

func NewYamlItemPrinter(w io.Writer) *DocItemPrinter {
	return &DocItemPrinter{
		w: w,
		exporterType: docExporterYaml,
	}
}

func (p *DocItemPrinter) Flush() {
	p.lock.Lock()

	var d []byte
	var err error
	if p.exporterType == docExporterJson {
		d, err = json.MarshalIndent(p.export, "", "  ")
		if err != nil {
			ReasonError.Printf("Couldn't produce JSON output: %v\n", err)
			return
		}
	} else if p.exporterType == docExporterYaml {
		d, err = yaml.Marshal(p.export)
		if err != nil {
			ReasonError.Printf("Couldn't produce JSON output: %v\n", err)
			return
		}
	}

	fmt.Fprintln(p.w, string(d))
	p.lock.Unlock()
}

func (p *DocItemPrinter) PrintItem(i *Item) {
	item := exportableItem{
		Region: i.Region.Name,
		Type: i.Type,
		State: i.State.String(),
	}

	rString, ok := i.Resource.(resources.LegacyStringer)
	if ok {
		item.Name = rString.String()
	}

	rProp, ok := i.Resource.(resources.ResourcePropertyGetter)
	if ok {
		item.Properties = rProp.Properties()
	}

	p.lock.Lock()
	p.export.Items = append(p.export.Items, item)
	p.lock.Unlock()
}