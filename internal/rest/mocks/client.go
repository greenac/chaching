package mocks

import (
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/rest/models"
	"net/http"
)

var _ models.IHttpClient = (*HttpClientMock)(nil)

type HttpClientMock struct {
	DoResponse http.Response
	DoError    error
}

func (c *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	return &c.DoResponse, c.DoError
}

var _ models.IClient = (*ClientMock)(nil)

type ClientMock struct {
	GetResponse      models.Response
	GetError         *genErr.GenError
	PostBodyResponse models.Response
	PostBodyError    *genErr.GenError
	PostUrlResponse  models.Response
	PostUrlError     *genErr.GenError
}

func (cm *ClientMock) Get(url string, headers *models.Headers, params models.UrlParams) (models.Response, genErr.IGenError) {
	return cm.GetResponse, cm.GetError
}

func (cm *ClientMock) PostBody(url string, headers *models.Headers, body []byte) (models.Response, genErr.IGenError) {
	return cm.PostBodyResponse, cm.PostBodyError
}

func (cm *ClientMock) PostUrl(url string, headers *models.Headers, params models.UrlParams) (models.Response, genErr.IGenError) {
	return cm.PostUrlResponse, cm.PostUrlError
}
