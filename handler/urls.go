package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ShaghayeghFathi/http-monitoring-service/model"
	"github.com/labstack/echo"
)

type urlResponse struct {
	ID           int       `json:"id"`
	URL          string    `json:"url"`
	UserID       uint      `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	Threshold    int       `json:"threshold"`
	FailedTimes  int       `json:"failed_times"`
	SuccessTimes int       `json:"success_times"`
}
type urlCreateRequest struct {
	Address      string `json:"address"`
	Threshold    int    `json:"threshold"`
	FailedTimes  int    `json:"failed_times"`
	SuccessTimes int    `json:"success_times"`
}

type urlListResponse struct {
	URLs     []*urlResponse `json:"urls"`
	UrlCount int            `json:"url_count"`
}

type urlStatusRequest struct {
	FromTime int64 `json:"from_time"`
	ToTime   int64 `json:"to_time"`
}

type requestResponse struct {
	ResultCode int       `json:"result_code"`
	CreatedAt  time.Time `json:"created_at"`
}

type requestListResponse struct {
	URL           string             `json:"url"`
	RequestsCount int                `json:"requests_count"`
	Requests      []*requestResponse `json:"requests"`
}

func newURLResponse(url *model.Url) *urlResponse {
	return &urlResponse{
		ID:           int(url.ID),
		URL:          url.Address,
		UserID:       url.UserId,
		Threshold:    url.Threshold,
		FailedTimes:  url.FailedTimes,
		SuccessTimes: url.SuccessTimes,
		CreatedAt: url.Model.CreatedAt,
	}
}

func newURLListResponse(list []model.Url) *urlListResponse {
	urls := make([]*urlResponse, 0)
	for i := range list {
		urls = append(urls, newURLResponse(&list[i]))
	}
	return &urlListResponse{
		URLs:     urls,
		UrlCount: len(list),
	}
}

func newRequestListResponse(reqs []model.Request, url string) *requestListResponse {
	resp := new(requestListResponse)
	resp.Requests = make([]*requestResponse, len(reqs))
	for i := range reqs {
		resp.Requests[i] = &requestResponse{
			ResultCode: reqs[i].Result,
			CreatedAt:  reqs[i].CreatedAt,
		}
	}
	resp.URL = url
	resp.RequestsCount = len(reqs)
	return resp
}

func bindToUrlCreateRequest(c echo.Context) (*urlCreateRequest, error) {
	request := &urlCreateRequest{}
	if err := c.Bind(request); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "error binding url create request" , err)
	}
	return request, nil
}

func bindToUrlStatusRequest(c echo.Context) (*urlStatusRequest, error) {
	req := &urlStatusRequest{}
	if err := c.Bind(req); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "error parsing url status request" , err)
	}
	return req, nil
}

func (h *Handler) FetchURLs(c echo.Context) error {
	userID := extractID(c)
	urls, err := h.dm.GetUrlsByUser(userID)
	if err != nil {
		return 	echo.NewHTTPError(http.StatusBadRequest, "Error retrieving urls from database, check token" , err)

	}
	resp := newURLListResponse(urls)
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateURL(c echo.Context) error {
	userID := extractID(c)

	req, err := bindToUrlCreateRequest(c)
	if err != nil {
		return err
	}

	url := &model.Url{
		UserId:       userID,
		Address:      req.Address,
		Threshold:    req.Threshold,
		FailedTimes:  req.FailedTimes,
		SuccessTimes: req.SuccessTimes,
	}

	if err := h.dm.AddUrl(url); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error adding url to database", err)
	}

	h.sch.Mnt.AddURL([]*model.Url{url})
	return c.JSON(http.StatusCreated, "URL created successfully")
}

func (h *Handler) GetURLStats(c echo.Context) error {
	userID := extractID(c)
	urlID, err := strconv.Atoi(c.Param("urlID"))
	if err != nil {
		return 	echo.NewHTTPError(http.StatusBadRequest, "Invalid path parameter" , err)

	}

	req, err := bindToUrlStatusRequest(c)
	if err != nil {
		return err
	}
	var url *model.Url
	if req.FromTime != 0 {
		if req.ToTime == 0 {
			req.ToTime = time.Now().Unix()
		}
		from := time.Unix(req.FromTime, 0)
		to := time.Unix(req.ToTime, 0)
		url, err = h.dm.GetUserRequestsInPeriod(uint(urlID), from, to)
	} else {
		url, err = h.dm.GetUrlById(uint(urlID))
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error retrieving url stats, invalid url id", err)
	}
	if url.UserId != userID {
		return echo.NewHTTPError(http.StatusUnauthorized, "operation not permitted")
	}
	return c.JSON(http.StatusOK, newRequestListResponse(url.Requests, url.Address))
}
