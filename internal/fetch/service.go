package fetch

import (
	"fmt"
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/rest"
	"github.com/greenac/chaching/internal/utils"
	"net/http"
)

func NewFetchService() IFetchService {
	return &FetchService{}
}

type IFetchService interface {
	Fetch(params rest.UrlParams) ([]byte, *genErr.GenError)
}

var _ IFetchService = (*FetchService)(nil)

type FetchService struct {
	Url        string
	RestClient rest.IClient
}

func (fc *FetchService) Fetch(params rest.UrlParams) ([]byte, *genErr.GenError) {
	resp, err := fc.RestClient.Get(fc.Url, nil, params)
	if err != nil {
		return []byte{}, err.AddMsg("FetchService:Fetch:failed to get")
	}

	if !utils.SliceContains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, resp.StatusCode) {
		ge := genErr.GenError{}
		return []byte{}, ge.AddMsg(fmt.Sprintf("FetchService:Fetch failed with code: %d and status %s for url: %s", resp.StatusCode, resp.Status, fc.Url))
	}

	return resp.Body, nil
}
