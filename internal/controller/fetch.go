package controller

import (
	"context"
	"fmt"
	"github.com/greenac/chaching/internal/database/models"
	"github.com/greenac/chaching/internal/database/service"
	genErr "github.com/greenac/chaching/internal/error"
	model "github.com/greenac/chaching/internal/rest/polygon/models"
	"github.com/greenac/chaching/internal/service/fetch"
	"github.com/greenac/chaching/internal/service/logger"
	"strings"
	"sync"
	"time"
)

const SortDirection = model.PolygonAggregateSortDirectionAsc
const ConcurrencyCount = 5
const logCount = 200

type FetchParams struct {
	TimespanMultiplier int
	Limit              int
	Timespan           model.PolygonAggregateTimespan
}

type FetchTargetParams struct {
	FetchParams
	CompanyName string
	From        time.Time
	To          time.Time
}

type FetchTarget struct {
	Name string
}

type FetchTargetsRetVal struct {
	DataPoints []models.DataPoint
	Error      genErr.IGenError
}

type FetchController struct {
	Targets         []string
	StartOfDay      time.Time
	EndOfDay        time.Time
	StartDate       time.Time
	EndDate         time.Time
	PartitionValue  time.Duration
	FetchService    fetch.FetchService
	DatabaseService service.IDatabaseService
	Logger          logger.ILogger
	Unmarshaler     func(data []byte, v any) error
}

func (fc *FetchController) RunFetch(fp FetchParams) []genErr.IGenError {
	var genErrs []genErr.IGenError

	times := fc.partitionTimes()
	fc.Logger.Info(fmt.Sprintf("FetchController:RunFetch:fetch # of times: %d", len(times)))
	for i := 0; i < len(times)-1; i += ConcurrencyCount {
		wg := sync.WaitGroup{}

		for j := i; j < i+ConcurrencyCount && j < len(times)-1; j += 1 {
			if j%logCount == 0 {
				fc.Logger.Info(fmt.Sprintf("FetchController:RunFetch:fetching from: %s to: %s. Count: %d to %d", times[j].Format(time.RFC3339), times[j+1].Format(time.RFC3339), j, j+ConcurrencyCount))
			}

			wg.Add(1)

			go func(from time.Time, to time.Time) {
				defer wg.Done()
				dps, errs := fc.FetchGroup(fp, from, to)
				if errs != nil {
					genErrs = append(genErrs, errs...)
				}

				gErrs := fc.DatabaseService.SaveDataPoints(context.Background(), dps)
				if gErrs != nil {
					genErrs = append(genErrs, *gErrs...)
				}
			}(times[j], times[j+1])
		}

		wg.Wait()
	}

	return genErrs
}

func (fc *FetchController) FetchGroup(fp FetchParams, from time.Time, to time.Time) ([]models.DataPoint, []genErr.IGenError) {
	var genErrors []genErr.IGenError
	var dataPts []models.DataPoint

	wg := sync.WaitGroup{}
	c := make(chan FetchTargetsRetVal)

	go func() {
		wg.Wait()
		close(c)
	}()

	for _, name := range fc.Targets {
		wg.Add(1)
		go func(n string) {
			defer func() {
				if r := recover(); r != nil {
					fc.Logger.Error(fmt.Sprintf("FetchController:FetchGroup:panic recovered for name: %s, fetch params: %+v, from: %s, to: %s, panic: %+v", n, fp, from.Format(time.RFC3339), to.Format(time.RFC3339), r))
				}
			}()

			defer wg.Done()
			fc.FetchTargets(FetchTargetParams{FetchParams: fp, CompanyName: n, From: from, To: to}, c)
		}(name)
	}

	for rv := range c {
		if rv.Error != nil {
			fc.Logger.Error(rv.Error.Error())
		} else {
			dataPts = append(dataPts, rv.DataPoints...)
		}
	}

	return dataPts, genErrors
}

func (fc *FetchController) FetchTargets(fp FetchTargetParams, c chan FetchTargetsRetVal) {
	rps := model.PolygonAggregateRequestParams{
		CompanyName:   fp.CompanyName,
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
			Error: &genErr.GenError{Messages: []string{"FetchController:FetchTargets:failed to unmarshall with err: " + err.Error() + " for: " + fp.CompanyName}},
		}
		return
	}

	if strings.ToLower(pr.Status) != "ok" {
		c <- FetchTargetsRetVal{
			Error: &genErr.GenError{Messages: []string{"FetchController:FetchTargets:failed for: " + fp.CompanyName + " with status: " + pr.Status + " at time: " + fp.From.Format(time.RFC3339)}},
		}
		return
	}

	dps := make([]models.DataPoint, len(pr.DataPoints))
	for i, dp := range pr.DataPoints {
		dps[i] = models.DataPoint{CompanyName: fp.CompanyName, PolygonDataPoint: dp}
	}

	c <- FetchTargetsRetVal{DataPoints: dps}
}

func (fc *FetchController) partitionTimes() []time.Time {
	var times []time.Time
	t := fc.StartDate
	startOfDay := fc.StartDate
	endOfDay := fc.EndOfDay
	for t.Before(fc.EndDate) || t.Equal(fc.EndDate) {
		if t.After(endOfDay) || t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
			startOfDay = startOfDay.Add(24 * time.Hour)
			endOfDay = endOfDay.Add(24 * time.Hour)
			t = startOfDay
			continue
		}

		times = append(times, t)
		t = t.Add(fc.PartitionValue)
	}

	return times
}
