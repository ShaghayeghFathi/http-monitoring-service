package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordValidation(t *testing.T) {
	user, err := NewUser("test", "Hello")
	assert.NoError(t, err, "Error creating user")
	assert.True(t, user.ValidatePassword("Hello"), "Error validating password ")
}

func TestUserCreation(t *testing.T) {
	_, err := NewUser("", "")
	assert.Error(t, err, "error creating user, username and password are empty")
}

func TestHashPassword(t *testing.T) {
	_, err := HashPassword("")
	assert.Error(t, err, "error hashing error")
}

func TestURLCreation(t *testing.T) {
	url, err := NewURL(0, "google.com", 10)
	assert.NoError(t, err, "error creating url")
	assert.Equal(t, url.Address, "http://google.com")
	_, err = NewURL(0, "ht://foo.bar", 10)
	assert.Error(t, err, "error validating url")
}

func TestURLSendRequest(t *testing.T) {
	url, _ := NewURL(0, "127.0.0.1:3000", 5)
	_, err := url.SendRequest()
	assert.Error(t, err)
	url.Address = "http://google.com"
	req, err := url.SendRequest()
	assert.NoError(t, err)
	assert.Equal(t, req.Result/100, 2)
}

func TestAlarmTrigger(t *testing.T) {
	url, _ := NewURL(0, "google.com", 2)
	url.FailedTimes = 2
	assert.True(t, url.ShouldTriggerAlarm(), "error triggering alarm")
}
