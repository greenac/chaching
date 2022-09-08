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
const ConcurrencyCount = 5

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
	End          time.Time
	FetchService fetch.FetchService
	Logger       logger.ILogger
	Unmarshaler  func(data []byte, v any) error
}

func (fc *FetchController) RunFetch(fp FetchParams) ([][]model.PolygonDataPoint, []genErr.IGenError) {
	genErrors := []genErr.IGenError{}
	dataPts := [][]model.PolygonDataPoint{}
	wg := sync.WaitGroup{}

	start := fc.Start
	partitions := int(fc.End.Sub(fc.Start).Hours())
	for i := 0; i < partitions; i += ConcurrencyCount {
		for j := i; j < i+ConcurrencyCount && j < partitions; j += 1 {
			wg.Add(1)

			go func(from time.Time) {
				defer wg.Done()
				dps, errs := fc.FetchGroup(fp, from, from.Add(time.Hour))
				if errs != nil {
					genErrors = append(genErrors, errs...)
				}
				dataPts = append(dataPts, dps...)
			}(start)

			start = start.Add(time.Hour)
		}
	}

	fmt.Println("FetchController:RunFetch:starting wait group")

	wg.Wait()

	fmt.Println("FetchController:RunFetch:leaving wait group")

	return dataPts, genErrors
}

func (fc *FetchController) FetchGroup(fp FetchParams, from time.Time, to time.Time) ([][]model.PolygonDataPoint, []genErr.IGenError) {
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
			fmt.Println("fetching for:", name, "from:", from.Format(time.RFC3339), "to:", to.Format(time.RFC3339))
			fc.FetchTargets(FetchTargetParams{FetchParams: fp, Name: name, From: from, To: to}, c)
		}(t)
	}

	for rv := range c {
		if rv.Error != nil {
			fc.Logger.Error().Msg(rv.Error.Error())
		} else {
			dataPts = append(dataPts, rv.DataPoints)
		}
	}

	//for _, dps := range dataPts {
	//	for _, dp := range dps {
	//		fmt.Printf("%+v\n", dp)
	//	}
	//}

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
