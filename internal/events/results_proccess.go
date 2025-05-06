package events

import (
	"fmt"
	"skiers/internal/config"
	"skiers/pkg/logger"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	resultString = "[%s] %d %v %v %s"
)

type CompetitorEventData struct {
	mu   *sync.Mutex
	data map[int]*EventDataById
}

type EventDataById struct {
	Data map[int][]*Event
}

type LapData struct {
	Time  time.Duration
	Speed float64
}

type CompetitorLapData struct {
	Time            time.Duration
	CompetitorId    int
	Status          string
	LapTimes        []*LapData
	PenaltyLapTimes []*LapData
	LapShots        []int
}

func (cld *CompetitorLapData) PrintData() {
	status := formatStatusToString(cld.Status, cld.Time)

	lapData := formatLapTimeToString(cld.LapTimes)
	penaltyLapData := formatLapTimeToString(cld.PenaltyLapTimes)

	shotNumbers := formatNumberOfShotsToString(cld.LapShots)

	fmt.Println(fmt.Sprintf(resultString, status, cld.CompetitorId, lapData, penaltyLapData, shotNumbers))
}

func (cld *CompetitorLapData) getStatus(events []*Event, cfg *config.Config) {
	switch len(events) - 1 {
	case cfg.Laps:
		var totalTime time.Duration

		for _, lapData := range cld.LapTimes {
			totalTime += lapData.Time
		}

		cld.Time = totalTime
		cld.Status = ""
	case 0:
		cld.Time = 0
		cld.Status = "NotStarted"
	default:
		cld.Time = 0
		cld.Status = "NotFinished"
	}
}

func (cld *CompetitorLapData) getTime(eventsDataById *EventDataById, events []*Event, cfg *config.Config) {
	cld.LapTimes = make([]*LapData, len(eventsDataById.Data[10]))
	cld.PenaltyLapTimes = make([]*LapData, len(eventsDataById.Data[10]))

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	var penaltyCounter int

	for i := 0; i < len(events)-1; i++ {
		lapData := LapData{}
		penaltyLapData := LapData{}

		if i == 0 {
			startTime, _ := time.Parse("15:04:05.000", events[i].ExtraParams)

			lapData.Time = events[i+1].Time.Sub(startTime)
		} else {
			lapData.Time = events[i+1].Time.Sub(events[i].Time)

		}

		lapData.Speed = float64(cfg.LapLen) / lapData.Time.Seconds()

		if penaltyCounter < len(eventsDataById.Data[8]) {
			if events[i+1].Time.After(eventsDataById.Data[8][penaltyCounter].Time) {
				penaltyLapData.Time = eventsDataById.Data[9][penaltyCounter].Time.Sub(eventsDataById.Data[8][penaltyCounter].Time)
				penaltyLapData.Speed = float64(cfg.PenaltyLen) / penaltyLapData.Time.Seconds()

				penaltyCounter++
			}
		} else {
			penaltyLapData.Time = 0
			penaltyLapData.Speed = 0
		}

		cld.LapTimes[i] = &lapData
		cld.PenaltyLapTimes[i] = &penaltyLapData
	}
}

func (cld *CompetitorLapData) getNumberOfShots(eventsDataById *EventDataById, events []*Event) {
	cld.LapShots = make([]int, len(events)-1)

	shotEvents := make([]*Event, len(eventsDataById.Data[6]))
	copy(shotEvents, eventsDataById.Data[6])

	for i, event := range events[1:] {
		var lapShots int
		var shotEventsIdx int

		for idx, v := range shotEvents {
			if v.Time.Before(event.Time) {
				lapShots++
				shotEventsIdx = idx
			}
		}

		cld.LapShots[i] = lapShots
		shotEvents = shotEvents[shotEventsIdx+1:]
	}
}

func (ced *CompetitorEventData) Set(competitorId int, eventsById *EventDataById) {
	ced.mu.Lock()
	defer ced.mu.Unlock()

	ced.data[competitorId] = eventsById
}

func (ced *CompetitorEventData) Get() map[int]*EventDataById {
	return ced.data
}

func getCompetitorEventData(eventData *EventData) *CompetitorEventData {
	logger.GetLogger().Info("Start grouping events by competitors ID")

	ced := &CompetitorEventData{mu: &sync.Mutex{}, data: make(map[int]*EventDataById)}

	wg := sync.WaitGroup{}

	for _, competitorId := range eventData.GetCompetitorIds() {
		wg.Add(1)

		go func() {
			defer wg.Done()

			competitorEvents, _ := eventData.Get(competitorId)

			eventsById := getCompetitorEventDataById(competitorId, competitorEvents)
			ced.Set(competitorId, eventsById)
		}()
	}

	wg.Wait()

	return ced
}

func getCompetitorEventDataById(competitorId int, competitorEvents []*Event) *EventDataById {
	logger.GetLogger().Infof("Start grouping competitior(%d) events by event ID", competitorId)

	eventsById := make(map[int][]*Event)

	for _, event := range competitorEvents {
		eventsById[event.Type] = append(eventsById[event.Type], event)
	}

	return &EventDataById{Data: eventsById}
}

func CreateResultTable(eventData *EventData, cfg *config.Config) {
	competitorEventData := getCompetitorEventData(eventData)

	logger.GetLogger().Info("Results Table\n")

	for competitorId, eventsDataById := range competitorEventData.Get() {
		events := make([]*Event, len(eventsDataById.Data[10])+1)

		startEvent, _ := eventsDataById.Data[2]
		lapEvents, _ := eventsDataById.Data[10]

		events[0] = startEvent[0]
		for idx, event := range lapEvents {
			events[idx+1] = event
		}

		cld := CompetitorLapData{CompetitorId: competitorId}

		cld.getTime(eventsDataById, events, cfg)
		cld.getStatus(events, cfg)
		cld.getNumberOfShots(eventsDataById, events)
		cld.PrintData()
	}
}

func formatDurationToString(format string, d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := (d.Nanoseconds() / 1e6) % 1000

	return fmt.Sprintf(format, hours, minutes, seconds, milliseconds)
}

func formatStatusToString(status string, time time.Duration) string {
	var resultStatus string

	if status == "" {
		resultStatus = formatDurationToString("%02d:%02d:%02d.%03d", time)
	} else {
		resultStatus = status
	}

	return resultStatus
}

func formatLapTimeToString(ld []*LapData) string {
	var lapData strings.Builder

	lapData.WriteString("[")

	for i, lap := range ld {
		lapTime := formatDurationToString("%02d:%02d:%02d.%03d", lap.Time)

		if i < len(ld)-1 {
			lapData.WriteString(fmt.Sprintf("{%s, %.3f}, ", lapTime, lap.Speed))
		} else {
			lapData.WriteString(fmt.Sprintf("{%s, %.3f}", lapTime, lap.Speed))
		}

	}

	lapData.WriteString("]")

	return lapData.String()
}

func formatNumberOfShotsToString(shots []int) string {
	var shotData strings.Builder

	shotData.WriteString("[")

	for i, shot := range shots {
		if i < len(shots)-1 {
			shotData.WriteString(fmt.Sprintf("{%d/5}, ", shot))
		} else {
			shotData.WriteString(fmt.Sprintf("{%d/5}", shot))
		}
	}

	shotData.WriteString("]")

	return shotData.String()
}
