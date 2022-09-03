package fetch

import (
	"github.com/greenac/chaching/internal/rest"
)

func NewFetchService() IFetchService {
	return &FetchService{}
}

type IFetchService interface {
	Fetch()
}

var _ IFetchService = (*FetchService)(nil)

type FetchService struct {
	RestClient rest.IClient
}

func (fc *FetchService) Fetch() {

}
