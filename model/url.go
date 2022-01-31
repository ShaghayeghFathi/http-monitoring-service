package model

import (
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Url struct {
	gorm.Model
	UserId       uint   `gorm:"unique_index:index_addr_user"` 
	Address      string `gorm:"unique_index:index_addr_user"`
	Threshold    int
	FailedTimes  int
	SuccessTimes int
	Requests     []Request `gorm:"foreignkey:url_id"`
}

func (url *Url) Serialize() *SerializedUrl {
	serializedUrl := &SerializedUrl{
		UserId:       url.UserId,
		Address:      url.Address,
		Threshold:    url.Threshold,
		FailedTimes:  url.FailedTimes,
		SuccessTimes: url.SuccessTimes,
		CreatedAt:    url.Model.CreatedAt,
	}
	if url.Requests != nil {
		serializedUrl.Requests = make([]SerializedRequest, 0, len(url.Requests))
		for _, request := range url.Requests {
			serializedUrl.Requests = append(serializedUrl.Requests, *request.Serialize())
		}
	}
	return serializedUrl
}

type SerializedUrl struct {
	CreatedAt    time.Time
	UserId       uint
	Address      string
	Threshold    int
	FailedTimes  int
	SuccessTimes int
	Requests     []SerializedRequest
}

type Request struct {
	gorm.Model
	UrlId  uint
	Result int
}

type SerializedRequest struct {
	CreatedAt time.Time
	UrlId     uint
	Result    int
}

func (request *Request) Serialize() *SerializedRequest {
	return &SerializedRequest{
		CreatedAt: request.CreatedAt,
		UrlId:     request.UrlId,
		Result:    request.Result,
	}
}

func NewURL(userID uint, address string, threshold int) (*Url, error) {
	url := new(Url)
	url.UserId = userID
	url.Threshold = threshold
	url.FailedTimes = 0
	url.SuccessTimes = 0

	if !strings.HasPrefix("http://", address) {
		address = "http://" + address
	}
	url.Address = address
	return url, nil
}

func (url *Url) ShouldTriggerAlarm() bool {
	return url.FailedTimes >= url.Threshold
}

func (url *Url) SendRequest() (*Request, error) {
	resp, err := http.Get(url.Address)
	req := new(Request)
	req.UrlId = url.ID
	if err != nil {
		return req, err
	}
	req.Result = resp.StatusCode
	return req, nil
}
