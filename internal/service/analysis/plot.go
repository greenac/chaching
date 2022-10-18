package analysis

import (
	"github.com/greenac/chaching/internal/database/models"
	genErr "github.com/greenac/chaching/internal/error"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"sort"
	"time"
)

type IPlotService interface {
	Plot(point []models.DataPoint) error
}

type PlotService struct {
}

func (ps *PlotService) Plot(points []models.DataPoint) genErr.IGenError {
	if len(points) == 0 {
		ge := genErr.GenError{}
		return ge.AddMsg("PlotService:Plot:points are empty")
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].StartTime < points[j].StartTime
	})

	firstPnt := points[0]
	//lastPnt := points[len(points) - 1]

	p := plot.New()

	plotter.DefaultLineStyle.Width = vg.Points(1)
	plotter.DefaultGlyphStyle.Radius = vg.Points(2)

	//p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
	//	{Value: 0, Label: "0"}, {Value: 0.25, Label: ""}, {Value: 0.5, Label: "0.5"}, {Value: 0.75, Label: ""}, {Value: 1, Label: "1"},
	//})

	xTicks := []plot.Tick{}
	xys := make(plotter.XYs, len(points))
	for i, pt := range points {
		d := time.Unix(pt.StartTime/1000, 0)
		x := float64((pt.StartTime - firstPnt.StartTime) / 10000)
		xys[i] = plotter.XY{X: x, Y: pt.HighestPrice}
		if d.Minute() == 0 || d.Minute() == 30 {
			xTicks = append(xTicks, plot.Tick{Value: x, Label: d.Format("15:04")})
		}
	}

	p.X.Tick.Marker = plot.ConstantTicks(xTicks)

	line, err := plotter.NewLine(xys)
	if err != nil {
		ge := genErr.GenError{}
		return ge.AddMsg("PlotService:Plot:failed to make new line").AddMsg(err.Error())
	}

	scatter, err := plotter.NewScatter(xys)
	if err != nil {
		ge := genErr.GenError{}
		return ge.AddMsg("PlotService:Plot:failed to make new scatter plot").AddMsg(err.Error())
	}
	p.Add(line, scatter)

	err = p.Save(2000, 1000, "plots/"+time.Now().Format(time.RFC3339)+".png")
	if err != nil {
		ge := genErr.GenError{}
		return ge.AddMsg("PlotService:Plot:failed to save plot").AddMsg(err.Error())
	}

	return nil
}
