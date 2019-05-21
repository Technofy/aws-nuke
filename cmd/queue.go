package cmd

import (
	"fmt"

	"github.com/rebuy-de/aws-nuke/resources"
)

type ItemState int

func (s ItemState) String() string {
	switch s {
	case ItemStateNew:
		return "new"
	case ItemStatePending:
		return "pending"
	case ItemStateWaiting:
		return "waiting"
	case ItemStateFailed:
		return "failed"
	case ItemStateFiltered:
		return "filtered"
	case ItemStateFinished:
		return "finished"
	default:
		return "unknown"
	}
}

// States of Items based on the latest request to AWS.
const (
	ItemStateNew ItemState = iota
	ItemStatePending
	ItemStateWaiting
	ItemStateFailed
	ItemStateFiltered
	ItemStateFinished
)

// An Item describes an actual AWS resource entity with the current state and
// some metadata.
type Item struct {
	Resource resources.Resource

	State  ItemState
	Reason string

	Region Region
	Type   string
	Printer ItemPrinter
}

func (i *Item) Print() {
	i.Printer.PrintItem(i)
}

// List gets all resource items of the same resource type like the Item.
func (i *Item) List() ([]resources.Resource, error) {
	listers := resources.GetListers()
	return listers[i.Type](i.Region.Session)
}

func (i *Item) GetProperty(key string) (string, error) {
	if key == "" {
		stringer, ok := i.Resource.(resources.LegacyStringer)
		if !ok {
			return "", fmt.Errorf("%T does not support legacy IDs", i.Resource)
		}
		return stringer.String(), nil
	}

	getter, ok := i.Resource.(resources.ResourcePropertyGetter)
	if !ok {
		return "", fmt.Errorf("%T does not support custom properties", i.Resource)
	}

	return getter.Properties().Get(key), nil
}

func (i *Item) Equals(o resources.Resource) bool {
	iType := fmt.Sprintf("%T", i.Resource)
	oType := fmt.Sprintf("%T", o)
	if iType != oType {
		return false
	}

	iStringer, iOK := i.Resource.(resources.LegacyStringer)
	oStringer, oOK := o.(resources.LegacyStringer)
	if iOK != oOK {
		return false
	}
	if iOK && oOK {
		return iStringer.String() == oStringer.String()
	}

	iGetter, iOK := i.Resource.(resources.ResourcePropertyGetter)
	oGetter, oOK := o.(resources.ResourcePropertyGetter)
	if iOK != oOK {
		return false
	}
	if iOK && oOK {
		return iGetter.Properties().Equals(oGetter.Properties())
	}

	return false
}

type Queue []*Item

func (q Queue) CountTotal() int {
	return len(q)
}

func (q Queue) Count(states ...ItemState) int {
	count := 0
	for _, item := range q {
		for _, state := range states {
			if item.State == state {
				count = count + 1
				break
			}
		}
	}
	return count
}
