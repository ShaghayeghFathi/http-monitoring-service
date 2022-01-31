package db_manager

import (
	"time"

	"github.com/ShaghayeghFathi/http-monitoring-service/model"

	"github.com/jinzhu/gorm"
)

type DbManager struct {
	db *gorm.DB
}

func NewDbManager(database *gorm.DB) *DbManager {
	return &DbManager{db: database}
}


func (s *DbManager) GetUserByUserName(username string) (*model.User, error) {
	user := new(model.User)
	if err := s.db.Preload("Urls").Preload("Urls.Requests").First(user, model.User{Username: username}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *DbManager) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *DbManager) AddUser(user *model.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return  err
	}
	return nil
}


func (s *DbManager) AddUrl(url *model.Url) error {
	return s.db.Create(url).Error
}

func (s *DbManager) GetUrlById(id uint) (*model.Url, error) {
	url := new(model.Url)
	if err := s.db.First(url, id).Error; err != nil {
		return nil, err
	}
	requests := make([]model.Request, 0)
	s.db.Model(url).Association("Requests").Find(&requests)
	url.Requests = requests
	return url, nil
}
func (s *DbManager) GetUrlsByUser(userID uint) ([]model.Url, error) {
	var urls []model.Url
	if err := s.db.Model(&model.Url{}).Where("user_id = ?", userID).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}

func (s *DbManager) UpdateUrl(url *model.Url) error {
	return s.db.Model(url).Update(url).Error
}

func (s *DbManager) GetAllUrls() ([]*model.Url, error) {
	var urls []*model.Url
	if err := s.db.Model(&model.Url{}).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}

func (s *DbManager) IncrementFailed(url *model.Url) error {
	url.FailedTimes += 1
	return s.UpdateUrl(url)
}

func (s *DbManager) IncrementSuccess(url *model.Url) error {
	url.SuccessTimes += 1
	return s.UpdateUrl(url)
}

func (s *DbManager) AddRequest(req *model.Request) error {
	return s.db.Create(req).Error
}
func (s *DbManager) GetRequestsByUrl(urlID uint) ([]model.Request, error) {
	var requests []model.Request
	if err := s.db.Model(&model.Request{UrlId: urlID}).Where("url_id == ?", urlID).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *DbManager) GetUserRequestsInPeriod(urlID uint, from, to time.Time) (*model.Url, error) {
	url := &model.Url{}
	url.ID = urlID
	if err := s.db.Model(url).Preload("Requests", "created_at >= ? and created_at <= ?", from, to).First(url).Error; err != nil {
		return nil, err
	}
	return url, nil
}

func (s *DbManager) FetchAlerts(userID uint) ([]*model.Url, error) {
	var urls []*model.Url
	if err := s.db.Model(&model.Url{}).Where("user_id = ? and failed_times >= threshold", userID).Find(&urls).Error; err != nil {
		return nil, err
	}
	for _, url := range urls {
		requests := make([]model.Request, 0)
		s.db.Model(url).Association("Requests").Find(&requests)
		url.Requests = requests
	}
	return urls, nil
}
