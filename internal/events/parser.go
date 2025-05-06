package events

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"skiers/pkg/logger"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	patternString = `^\[(\d{2}:\d{2}:\d{2}\.\d{3})\]\s+(\d+)\s+(\d+)(?:\s+(.+))?$`
)

type EventData struct {
	mu   *sync.Mutex
	data map[int][]*Event
}

func NewEventData() *EventData {
	return &EventData{mu: &sync.Mutex{}, data: make(map[int][]*Event)}
}

func (ed *EventData) Set(event *Event) {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	ed.data[event.CompetitorID] = append(ed.data[event.CompetitorID], event)
}

func (ed *EventData) Get(competitorID int) ([]*Event, error) {
	events, ok := ed.data[competitorID]
	if !ok {
		return nil, fmt.Errorf("invalid competitor ID: %d", competitorID)
	}

	return events, nil
}

func (ed *EventData) GetAllEvents() []*Event {
	var fullEvents []*Event

	for _, competitorEvents := range ed.data {
		fullEvents = append(fullEvents, competitorEvents...)
	}

	return fullEvents
}

func (ed *EventData) GetCompetitorIds() []int {
	keys := make([]int, 0, len(ed.data))
	for k := range ed.data {
		keys = append(keys, k)
	}
	return keys
}

type Event struct {
	Time         time.Time
	Type         int
	CompetitorID int
	ExtraParams  string
}

func ParseEventFile(path string) (*EventData, error) {
	g, ctx := errgroup.WithContext(context.Background())

	cpuNum := runtime.NumCPU()
	lines := readFile(ctx, path, cpuNum)

	eventData := NewEventData()

	for i := 0; i < cpuNum; i++ {
		g.Go(func() error {
			return worker(ctx, i, eventData, lines)
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return eventData, nil
}

func atoi(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func readFile(ctx context.Context, path string, cpuNum int) chan string {
	lines := make(chan string, cpuNum*2)

	go func() {
		defer close(lines)

		logger.GetLogger().Info("start reading event file")

		file, err := os.Open(path)
		if err != nil {
			logger.GetLogger().Errorf("error while opening event file: %v", err)
			return
		}

		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				logger.GetLogger().Errorf("error while closing event file: %v", err)
				return
			}
		}(file)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case lines <- scanner.Text():
			}
		}
	}()

	return lines
}

func parseLine(line string) (*Event, error) {
	matches := regexp.MustCompile(patternString).FindStringSubmatch(line)
	if matches == nil {
		return nil, errors.New("invalid event string")
	}

	eventTime, err := time.Parse("15:04:05.000", matches[1])
	if err != nil {
		return nil, fmt.Errorf("error while time parsing: %w", err)
	}

	return &Event{
		Time:         eventTime,
		Type:         atoi(matches[2]),
		CompetitorID: atoi(matches[3]),
		ExtraParams:  matches[4],
	}, nil
}

func worker(ctx context.Context, workerNumber int, eventData *EventData, lines chan string) error {
	logger.GetLogger().Infof("worker #%d starts parsing lines", workerNumber)

	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-lines:
			if !ok {
				return nil
			}

			event, err := parseLine(line)
			if err != nil {
				return err
			}

			eventData.Set(event)
		}
	}
}
