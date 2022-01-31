package service

import (
	"os"
	"testing"

	"github.com/ShaghayeghFathi/http-monitoring-service/db"
	"github.com/ShaghayeghFathi/http-monitoring-service/db_manager"
	"github.com/ShaghayeghFathi/http-monitoring-service/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

var mnt *Monitor
var st *db_manager.DbManager
var d *gorm.DB

func setupDB() {
	d = db.Setup("test-monitor.db")
	st = db_manager.NewDbManager(d)
	user, _ := model.NewUser("foo", "bar")
	_ = st.AddUser(user)

	mnt = NewMonitor(st)
}

func removeDB() {
	if d != nil {
		if err := d.Close(); err != nil {
			log.Error(err)
		}
	}
	if err := os.Remove("test-monitor.db"); err != nil {
		log.Error(err)
	}
}
func TestMonitor_Work(t *testing.T) {
	removeDB()
	setupDB()
	urls := []*model.Url{
		{UserId: 1, Address: "http://google.com", Threshold: 10, FailedTimes: 0},
		{UserId: 2, Address: "http://google.com", Threshold: 10, FailedTimes: 0},
	}
	err := st.AddUrl(urls[0])
	assert.NoError(t, err, "error creating first url")
	err = st.AddUrl(urls[1])
	assert.NoError(t, err, "error creating second url")
	mnt.AddURL(urls)
	mnt.Work()
	req, err := st.GetRequestsByUrl(urls[0].ID)
	assert.NoError(t, err, "error in GetRequestsByUrl")
	assert.Len(t, req, 1)
}
