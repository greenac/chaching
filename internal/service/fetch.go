package fetch

import (
	"fmt"
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/rest/models"
	"github.com/greenac/chaching/internal/utils"
	"net/http"
)

type IFetchData interface {
	Url() string
}

func NewFetchService(url string, rc models.IClient, joiner func(base string, add string) (result string, err *genErr.GenError)) *FetchService {
	return &FetchService{Url: url, RestClient: rc, PathJoiner: joiner}
}

type IFetchService interface {
	Fetch(params models.UrlParams) ([]byte, *genErr.GenError)
	FetchWithFetchData(fetchData IFetchData) ([]byte, *genErr.GenError)
}

var _ IFetchService = (*FetchService)(nil)

type FetchService struct {
	Url        string
	RestClient models.IClient
	PathJoiner func(base string, add string) (string, *genErr.GenError)
}

func (fc *FetchService) Fetch(params models.UrlParams) ([]byte, *genErr.GenError) {
	resp, err := fc.RestClient.Get(fc.Url, nil, params)
	if err != nil {
		return []byte{}, err.AddMsg("FetchService:Fetch:failed to get")
	}

	return fc.handleResponse(resp)
}

func (fc *FetchService) FetchWithFetchData(fetchData IFetchData) ([]byte, *genErr.GenError) {
	uri, ge := fc.PathJoiner(fc.Url, fetchData.Url())
	if ge != nil {
		return []byte{}, ge.AddMsg("FetchService:FetchWithFetchData failed to join urls: " + fc.Url + " and " + fetchData.Url())
	}

	resp, ge := fc.RestClient.Get(uri, nil, models.UrlParams{})
	if ge != nil {
		return []byte{}, ge.AddMsg("FetchService:Fetch:failed to get")
	}

	return fc.handleResponse(resp)
}

func (fc *FetchService) handleResponse(resp models.Response) ([]byte, *genErr.GenError) {
	if !utils.SliceContains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, resp.StatusCode) {
		ge := genErr.GenError{}
		return []byte{}, ge.AddMsg(fmt.Sprintf("FetchService:handleResponse failed with code: %d and status %s for url: %s", resp.StatusCode, resp.Status, fc.Url))
	}

	return resp.Body, nil
}
