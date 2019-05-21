package cmd

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/rebuy-de/aws-nuke/resources"
)

var (
	ReasonSkip            = *color.New(color.FgYellow)
	ReasonError           = *color.New(color.FgRed)
	ReasonRemoveTriggered = *color.New(color.FgGreen)
	ReasonWaitPending     = *color.New(color.FgBlue)
	ReasonSuccess         = *color.New(color.FgGreen)
)

var (
	ColorRegion             = *color.New(color.Bold)
	ColorResourceType       = *color.New()
	ColorResourceID         = *color.New(color.Bold)
	ColorResourceProperties = *color.New(color.Italic)
)

type LogItemPrinter struct {
	w io.Writer
}

func NewLogItemPrinter(w io.Writer) *LogItemPrinter {
	return &LogItemPrinter{
		w: w,
	}
}

func (p *LogItemPrinter) Flush() {

}

func (p *LogItemPrinter) PrintItem(i *Item) {
	switch i.State {
	case ItemStateNew:
		p.logItem(i.Region, i.Type, i.Resource, ReasonWaitPending, "would remove")
	case ItemStatePending:
		p.logItem(i.Region, i.Type, i.Resource, ReasonWaitPending, "triggered remove")
	case ItemStateWaiting:
		p.logItem(i.Region, i.Type, i.Resource, ReasonWaitPending, "waiting")
	case ItemStateFailed:
		p.logItem(i.Region, i.Type, i.Resource, ReasonError, "failed")
	case ItemStateFiltered:
		p.logItem(i.Region, i.Type, i.Resource, ReasonSkip, i.Reason)
	case ItemStateFinished:
		p.logItem(i.Region, i.Type, i.Resource, ReasonSuccess, "removed")
	}
}

func (p *LogItemPrinter) logItem(region Region, resourceType string, r resources.Resource, c color.Color, msg string) {
	ColorRegion.Fprintf(p.w, "%s", region.Name)
	fmt.Printf(" - ")
	ColorResourceType.Fprint(p.w, resourceType)
	fmt.Printf(" - ")

	rString, ok := r.(resources.LegacyStringer)
	if ok {
		ColorResourceID.Fprint(p.w, rString.String())
		fmt.Printf(" - ")
	}

	rProp, ok := r.(resources.ResourcePropertyGetter)
	if ok {
		ColorResourceProperties.Fprint(p.w, rProp.Properties())
		fmt.Printf(" - ")
	}

	c.Fprintf(p.w, "%s\n", msg)
}
