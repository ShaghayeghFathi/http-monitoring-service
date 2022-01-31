package db_manager

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ShaghayeghFathi/http-monitoring-service/db"
	"github.com/ShaghayeghFathi/http-monitoring-service/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

var database *gorm.DB
var st *DbManager
var usersList []*model.User
var urlsList []*model.Url

func TestMain(m *testing.M) {
	database = db.Setup("test.db")
	st = NewDbManager(database)

	setup()

	returnCode := m.Run()
	if err := database.Close(); err != nil {
		log.Error(err)
	}
	if err := os.Remove("test.db"); err != nil {
		log.Error(err)
	}

	os.Exit(returnCode)
}



func setup() {
	usersList = make([]*model.User, 2)
	usersList[0], _ = model.NewUser("user1", "123456")
	usersList[1], _ = model.NewUser("user2", "1234")

	urlsList = make([]*model.Url, 10)
	for i := range urlsList {
		urlsList[i] = new(model.Url)
		urlsList[i].UserId = usersList[0].ID
		urlsList[i].Address = fmt.Sprintf("www.%dsomething.com", i)
		urlsList[i].Threshold = 10
	}
}

func TestUsers(t *testing.T) {
	err := st.AddUser(usersList[0])
	assert.NoError(t, err, "error adding user to database")
	_ = st.AddUser(usersList[1])
	dbUser, err := st.GetUserByUserName("user1")
	assert.NoError(t, err, "error reading user from database")
	assert.Equal(t, dbUser.Username, "user1")
	_, err = st.GetUserByUserName("invalid-username")
	assert.Error(t, err)
	users, err := st.GetAllUsers()
	assert.NoError(t, err, "error reading all users from database")
	assert.Equal(t, 2, len(users))
	usersList[0], usersList[1] = &users[0], &users[1]
}

func TestUrls(t *testing.T) {

	for i := range urlsList {
		urlsList[i].UserId = usersList[0].ID
		err := st.AddUrl(urlsList[i])
		assert.NoError(t, err, "Error inserting url into database")
	}

	u, err := st.GetUrlById(1)
	assert.NoError(t, err, "Error reading url with id 1 from database")

	assert.Equal(t, u.Address, "www.0something.com", "Mismatch url in database")

	_, err = st.GetUrlById(1000)
	assert.Error(t, err)
	// Updating URL

	_, err = st.GetUrlsByUser(usersList[0].ID)
	assert.NoError(t, err)

	err = st.IncrementFailed(u)
	err = st.IncrementFailed(u)
	assert.NoError(t, err, "Error incrementing failed times")

	u, _ = st.GetUrlById(1)
	assert.Equal(t, 2, u.FailedTimes, "Increment failed_times didn't work")
}

func TestRequests(t *testing.T) {
	for i := range urlsList {
		st.AddUrl(urlsList[i])
		req := new(model.Request)
		req.Result = 300
		req.UrlId = urlsList[i/3].ID
		err := st.AddRequest(req)
		assert.NoError(t, err)
	}
	reqs, err := st.GetRequestsByUrl(urlsList[0].ID)
	assert.NoError(t, err, "Error retrieving requests from database")
	urlsByTime, err := st.GetUserRequestsInPeriod(urlsList[0].ID, time.Now().Add(-time.Minute*3), time.Now())
	assert.NoError(t, err)
	reqIDs := make([]uint, 0)
	for _, req := range reqs {
		reqIDs = append(reqIDs, req.ID)
	}
	assert.Equal(t, 3, len(reqs), "Mismatch between number of inserted and retrieved requests. %v %v", len(reqs), reqIDs)
	assert.Equal(t, 3, len(urlsByTime.Requests), "error getting urls filtered by time")
}
