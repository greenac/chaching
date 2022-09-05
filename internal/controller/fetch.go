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

type FetchController struct {
	Targets      []string
	Start        time.Time
	Delta        time.Duration
	FetchService fetch.FetchService
	Logger       logger.ILogger
	Unmarshaler  func(data []byte, v any) error
}

func (fc *FetchController) RunFetch(fp FetchParams) ([][]model.PolygonDataPoint, []*genErr.GenError) {
	genErrors := []*genErr.GenError{}
	dataPts := [][]model.PolygonDataPoint{}

	wg := sync.WaitGroup{}
	cErr := make(chan *genErr.GenError)
	cDataPoints := make(chan []model.PolygonDataPoint)

	go func() {
		wg.Wait()
		close(cErr)
		close(cDataPoints)
	}()

	for _, t := range fc.Targets {
		wg.Add(1)

		go func(name string) {
			defer wg.Done()
			fc.Logger.Info().Msg("FetchController::RunFetch fetching " + name)
			fc.FetchTargets(
				FetchTargetParams{FetchParams: fp, Name: name, From: fc.Start, To: fc.Start.Add(fc.Delta)},
				cErr,
				cDataPoints,
			)
		}(t)
	}

	for e := range cErr {
		if e != nil {
			genErrors = append(genErrors, e)
		}

		fc.Logger.Error().Msg("FetchController::RunFetch failed with error: " + e.Error())
	}

	for dps := range cDataPoints {
		dataPts = append(dataPts, dps)
		fc.Logger.Info().Msg(fmt.Sprintf("FetchController::RunFetch got data: %+v", dataPts))
	}

	return dataPts, genErrors
}

func (fc *FetchController) FetchTargets(fp FetchTargetParams, cErr chan *genErr.GenError, cDataPoint chan []model.PolygonDataPoint) {
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
		cErr <- ge
		return
	}

	pr := model.PolygonAggregateResponse{}
	err := fc.Unmarshaler(body, &pr)
	if err != nil {
		cErr <- &genErr.GenError{Messages: []string{"FetchController:FetchTargets failed to unmarshall with err: " + err.Error() + " for: " + fp.Name}}
		return
	}

	if strings.ToLower(pr.Status) != "ok" {
		cErr <- &genErr.GenError{Messages: []string{"FetchController:FetchTargets failed for: " + fp.Name + " with status: " + pr.Status}}
		return
	}

	cDataPoint <- pr.DataPoints
}
