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
	"github.com/greenac/chaching/internal/worker"
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

type FetchTaskResult struct {
	DataPoints []models.DataPoint
	Errors     *[]genErr.IGenError
}

func (fc *FetchController) RunFetch(fp FetchParams) []genErr.IGenError {
	msgChan := make(chan worker.Message[FetchTaskResult])

	times := fc.partitionTimes()
	fc.Logger.Info(fmt.Sprintf("FetchController:RunFetch:fetch # of times: %d", len(times)))

	wrkr := worker.NewWorker(10, msgChan)
	wrkr.Work()

	go func() {
		for i := 0; i < len(times)-1; i += 2 {
			go func(from time.Time, to time.Time) {
				task := func() FetchTaskResult {
					var genErrs []genErr.IGenError
					dps, errs := fc.FetchGroup(fp, from, to)
					if errs != nil {
						genErrs = append(genErrs, errs...)
					}

					gErrs := fc.DatabaseService.SaveDataPoints(context.Background(), dps)
					if gErrs != nil {
						genErrs = append(genErrs, *gErrs...)
					}

					return FetchTaskResult{DataPoints: dps, Errors: gErrs}
				}
				wrkr.AddTask(task)
			}(times[i], times[i+1])
		}
	}()

	var errors []genErr.IGenError
	for msg := range msgChan {
		if msg.Result.Errors != nil && len(*msg.Result.Errors) > 0 {
			errors = append(errors, *msg.Result.Errors...)
		}
	}

	return errors
}

func (fc *FetchController) FetchGroup(fp FetchParams, from time.Time, to time.Time) ([]models.DataPoint, []genErr.IGenError) {
	var genErrors []genErr.IGenError
	var dataPts []models.DataPoint

	c := make(chan FetchTargetsRetVal, len(fc.Targets))

	for _, name := range fc.Targets {
		go func(n string) {
			defer func() {
				if r := recover(); r != nil {
					fc.Logger.Error(fmt.Sprintf("FetchController:FetchGroup:panic recovered for name: %s, fetch params: %+v, from: %s, to: %s, panic: %+v", n, fp, from.Format(time.RFC3339), to.Format(time.RFC3339), r))
				}
			}()

			dps, gErr := fc.FetchTargets(FetchTargetParams{FetchParams: fp, CompanyName: n, From: from, To: to})
			c <- FetchTargetsRetVal{DataPoints: dps, Error: gErr}
		}(name)
	}

	for r := range c {
		dataPts = append(dataPts, r.DataPoints...)
		if r.Error != nil {
			genErrors = append(genErrors, r.Error)
		}
	}

	return dataPts, genErrors
}

func (fc *FetchController) FetchTargets(fp FetchTargetParams) ([]models.DataPoint, genErr.IGenError) {
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
		return []models.DataPoint{}, ge
	}

	pr := model.PolygonAggregateResponse{}
	err := fc.Unmarshaler(body, &pr)
	if err != nil {
		return []models.DataPoint{}, genErr.GenError{Messages: []string{"FetchController:FetchTargets:failed to unmarshall with err: " + err.Error() + " for: " + fp.CompanyName}}
	}

	if strings.ToLower(pr.Status) != "ok" {
		c <- FetchTargetsRetVal{
			Error: &genErr.GenError{Messages: []string{"FetchController:FetchTargets:failed for: " + fp.CompanyName + " with status: " + pr.Status + " at time: " + fp.From.Format(time.RFC3339)}},
		}
		return []models
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
