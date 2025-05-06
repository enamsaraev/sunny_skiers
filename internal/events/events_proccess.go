package events

import (
	"fmt"
	"sort"
	"strings"
)

type eventMapString struct {
	m map[int]string
}

func (ems *eventMapString) getString(e *Event) (string, bool) {
	v, ok := ems.m[e.Type]
	if !ok {
		return "", false
	}

	switch e.Type {
	case 1:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID), true
	case 2:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID, e.ExtraParams), true
	case 3:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID), true
	case 4:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID), true
	case 5:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID, e.ExtraParams), true
	case 6:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.ExtraParams, e.CompetitorID), true
	case 7:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID), true
	case 8:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID), true
	case 9:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID), true
	case 10:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID), true
	case 11:
		return fmt.Sprintf(v, e.Time.Format("15:04:05.000"), e.CompetitorID, e.ExtraParams), true
	default:
		return "", false
	}
}

func newEventMapString() *eventMapString {
	return &eventMapString{m: map[int]string{
		1:  "[%s] The competitor(%d) registered\n",
		2:  "[%s] The start time for the competitor(%d) was set by a draw to %s\n",
		3:  "[%s] The competitor(%d) is on the start line\n",
		4:  "[%s] The competitor(%d) has started\n",
		5:  "[%s] The competitor(%d) is on the firing range(%s)\n",
		6:  "[%s] The target(%s) has been hit by competitor(%d)\n",
		7:  "[%s] The competitor(%d) left the firing range\n",
		8:  "[%s] The competitor(%d) entered the penalty laps\n",
		9:  "[%s] The competitor(%d) left the penalty laps\n",
		10: "[%s] The competitor(%d) ended the main lap\n",
		11: "[%s] The competitor(%d) can`t continue: %s\n",
	}}
}

func LogCompetitorsData(eventData *EventData) error {
	sortedEvents := getSortedEventsByTime(eventData)

	ems := newEventMapString()

	var builder strings.Builder
	for _, e := range sortedEvents {
		eventString, ok := ems.getString(e)
		if !ok {
			fmt.Println("error", e)
		}
		builder.WriteString(eventString)
	}
	fmt.Println(builder.String())

	return nil
}

func getSortedEventsByTime(eventData *EventData) []*Event {
	fullEvents := eventData.GetAllEvents()

	sort.Slice(fullEvents, func(i, j int) bool {
		return fullEvents[i].Time.Before(fullEvents[j].Time)
	})

	return fullEvents
}
