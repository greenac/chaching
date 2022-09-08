package controller

import (
	"fmt"
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/logger"
	model "github.com/greenac/chaching/internal/rest/polygon/models"
	fetch "github.com/greenac/chaching/internal/service"
	"strings"
	"sync"
	"time"
)

const SortDirection = model.PolygonAggregateSortDirectionAsc

type FetchParams struct {
	TimespanMultiplier int
	Limit              int
	Timespan           model.PolygonAggregateTimespan
}

type FetchTargetParams struct {
	FetchParams
	Name string
	From time.Time
	To   time.Time
}

type FetchTarget struct {
	Name string
}

type FetchTargetsRetVal struct {
	DataPoints []model.PolygonDataPoint
	Error      genErr.IGenError
}

type FetchController struct {
	Targets      []string
	Start        time.Time
	Delta        time.Duration
	FetchService fetch.FetchService
	Logger       logger.ILogger
	Unmarshaler  func(data []byte, v any) error
}

func (fc *FetchController) RunFetch(fp FetchParams) ([][]model.PolygonDataPoint, []genErr.IGenError) {
	genErrors := []genErr.IGenError{}
	dataPts := [][]model.PolygonDataPoint{}

	wg := sync.WaitGroup{}
	c := make(chan FetchTargetsRetVal)

	go func() {
		wg.Wait()
		close(c)
	}()

	for _, t := range fc.Targets {
		wg.Add(1)

		go func(name string) {
			defer wg.Done()
			fc.FetchTargets(FetchTargetParams{FetchParams: fp, Name: name, From: fc.Start, To: fc.Start.Add(fc.Delta)}, c)
		}(t)
	}

	fc.Logger.Info().Msg("for loop before channel")

	for rv := range c {
		if rv.Error != nil {
			fc.Logger.Error().Msg(rv.Error.Error())
		} else {
			dataPts = append(dataPts, rv.DataPoints)
		}
	}

	for _, dps := range dataPts {
		for _, dp := range dps {
			fmt.Printf("%+v\n", dp)
		}
	}

	fc.Logger.Info().Msg("all done!")

	return dataPts, genErrors
}

func (fc *FetchController) FetchTargets(fp FetchTargetParams, c chan FetchTargetsRetVal) {
	rps := model.PolygonAggregateRequestParams{
		Name:          fp.Name,
		Multiplier:    fp.TimespanMultiplier,
		Timespan:      fp.Timespan,
		From:          fp.From,
		To:            fp.To,
		SortDirection: SortDirection,
		Limit:         fp.Limit,
	}

	body, ge := fc.FetchService.FetchWithFetchData(rps)
	if ge != nil {
		c <- FetchTargetsRetVal{Error: ge}
		return
	}

	pr := model.PolygonAggregateResponse{}
	err := fc.Unmarshaler(body, &pr)
	if err != nil {
		c <- FetchTargetsRetVal{
			Error: &genErr.GenError{Messages: []string{"FetchController:FetchTargets failed to unmarshall with err: " + err.Error() + " for: " + fp.Name}},
		}
		return
	}

	if strings.ToLower(pr.Status) != "ok" {
		c <- FetchTargetsRetVal{
			Error: &genErr.GenError{Messages: []string{"FetchController:FetchTargets failed for: " + fp.Name + " with status: " + pr.Status}},
		}
		return
	}

	c <- FetchTargetsRetVal{DataPoints: pr.DataPoints}
}
