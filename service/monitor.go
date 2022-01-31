package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ShaghayeghFathi/http-monitoring-service/db_manager"
	"github.com/ShaghayeghFathi/http-monitoring-service/model"
	"golang.org/x/sync/semaphore"
)

type Monitor struct {
	dm      *db_manager.DbManager
	urls       []*model.Url
	maxSemSize int
	sem        *semaphore.Weighted
}

func NewMonitor(dm *db_manager.DbManager) *Monitor {
	mnt := new(Monitor)
	mnt.dm = dm
	mnt.maxSemSize = 10
	mnt.sem = semaphore.NewWeighted(int64(10))
	mnt.getUrlFromDB()
	return mnt
}

func (mnt *Monitor) getUrlFromDB() error {
	urls, err := mnt.dm.GetAllUrls()
	if err != nil {
		return err
	}
	mnt.urls = urls
	return nil
}

func (mnt *Monitor) Work() {
	var wg sync.WaitGroup

	for urlIndex := range mnt.urls {
		url := mnt.urls[urlIndex]
		wg.Add(1)
		go func(urlIndex int) {
			if err := mnt.sem.Acquire(context.Background(), 1); err != nil {
				log.Fatal(err)
			}
			defer wg.Done()
			mnt.monitorURL(url)
			defer mnt.sem.Release(1)
		}(urlIndex)
	}
	wg.Wait()
}

func (mnt *Monitor) AddURL(urls []*model.Url) {
	mnt.urls = append(mnt.urls, urls...)
}

func (mnt *Monitor) monitorURL(url *model.Url) {

	req, err := url.SendRequest()
	if err != nil {
		fmt.Println(err, "could not make request")
		req = new(model.Request)
		req.UrlId = url.ID
		req.Result = http.StatusBadRequest
	}
	if err = mnt.dm.AddRequest(req); err != nil {
		fmt.Println(err, "could not save request to database")
	}

	if req.Result/100 == 2 {
		if err = mnt.dm.IncrementSuccess(url); err != nil {
			fmt.Println(err, "could not increment success times for url")
		}
	} else {
		if err = mnt.dm.IncrementFailed(url); err != nil {
			fmt.Println(err, "could not increment failed times for url")
		}
	}
}
