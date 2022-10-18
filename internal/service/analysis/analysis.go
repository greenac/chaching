package analysis

import (
	"github.com/greenac/chaching/internal/database/models"
	"github.com/greenac/chaching/internal/utils"
	"math"
	"time"
)

type SlopeChange struct {
	models.DataPoint
	Time            time.Time
	Slope           float64
	NormalizedSlope float64
}

func (sc *SlopeChange) NormalizeSlope(nf float64) {
	sc.NormalizedSlope = sc.Slope / nf
}

type slopeVal int

const (
	slopeValPos slopeVal = iota
	slopeValNeg
)

type IAnalysisService interface {
	FindSlopeChanges(points []models.DataPoint) []SlopeChange
	CalcSlopeNormalizationFactor(slopes []SlopeChange) float64
	SlopeAbsMidpoint(slopes []SlopeChange) float64
}

func NewAnalysisService() IAnalysisService {
	return &AnalysisService{}
}

var _ IAnalysisService = (*AnalysisService)(nil)

type AnalysisService struct{}

func (as *AnalysisService) FindSlopeChanges(points []models.DataPoint) []SlopeChange {
	if len(points) < 2 {
		return []SlopeChange{}
	}

	slopeChanges := []SlopeChange{}
	var slopeSign = slopeValPos

	for i := 0; i < len(points)-1; i += 1 {
		pt1 := points[i]
		pt2 := points[i+1]

		slope := as.calcSlope(pt1, pt2)
		ss := as.slopeSign(slope)
		if i == 0 {
			slopeSign = ss
		} else if ss != slopeSign {
			slopeSign = ss
			slopeChanges = append(slopeChanges, SlopeChange{Slope: slope, Time: pt1.Time(), DataPoint: pt2})
		}
	}

	return slopeChanges
}

func (as *AnalysisService) calcSlope(p1 models.DataPoint, p2 models.DataPoint) float64 {
	return (p2.YVal() - p1.YVal()) / (p2.XVal() - p1.XVal())
}

func (as *AnalysisService) slopeSign(slope float64) slopeVal {
	ss := slopeValPos
	if slope < 0 {
		ss = slopeValNeg
	}

	return ss
}

func (as *AnalysisService) CalcSlopeNormalizationFactor(slopes []SlopeChange) float64 {
	var max float64 = 0
	for _, s := range slopes {
		if math.Abs(s.Slope) > max {
			max = math.Abs(s.Slope)
		}
	}

	return max
}

func (as *AnalysisService) SlopeAbsMidpoint(slopes []SlopeChange) float64 {
	slps := make([]float64, len(slopes))
	for i, s := range slopes {
		slps[i] = math.Abs(s.NormalizedSlope)
	}

	return utils.MidPoint(slps)
}
