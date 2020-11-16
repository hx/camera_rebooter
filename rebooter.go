package main

import (
	"crypto/tls"
	"errors"
	"github.com/hx/logs"
	"net/http"
	"time"
)

type Config struct {
	Address    string `json:"address"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CameraName string `json:"camera_name"`
	Time       string `json:"time"`
	Location   string `json:"location"`
}

type Rebooter struct {
	*Config

	OnHour   int
	OnMinute int
	Location *time.Location

	Logger *logs.Logger

	cameraID   string
	authToken  string
	client     *http.Client
	nextReboot time.Time
}

func (r *Rebooter) init() (err error) {
	r.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 15 * time.Second,
	}
	r.authToken = ""
	r.cameraID = ""

	if err = r.Login(); err != nil {
		return
	}

	bs, err := r.Bootstrap()
	if err != nil {
		return
	}

	for _, cam := range bs.Cameras {
		if cam.Name == r.CameraName {
			r.Logger.Info("Camera '%s' has ID %s", cam.Name, cam.ID)
			r.cameraID = cam.ID
			break
		}
	}
	if r.cameraID == "" {
		return errors.New("camera name not found")
	}
	return
}

const TimeFormat = "3:04pm on Monday 2 January"

func (r *Rebooter) setNextScheduledTime() {
	now := time.Now().In(r.Location)
	year, month, day := now.Date()
	next := time.Date(year, month, day, r.OnHour, r.OnMinute, 0, 0, r.Location)
	if next.Before(now) {
		next = next.AddDate(0, 0, 1)
	}
	r.nextReboot = next
	r.Logger.Info("Next reboot will be at %s",
		next.Format(TimeFormat),
	)
}

func (r *Rebooter) Loop() (err error) {
	if err = r.init(); err != nil {
		return
	}
	r.setNextScheduledTime()
	for {
		for r.nextReboot.After(time.Now()) {
			d := r.nextReboot.Sub(time.Now())
			if d > time.Minute {
				d = time.Minute
			}
			<-time.After(d)
		}
		if err = r.Reboot(); err == nil {
			r.setNextScheduledTime()
		} else {
			r.Logger.Error(err.Error())
			r.nextReboot = time.Now().Add(time.Minute * 2)
			r.Logger.Info("Next retry will be at %s",
				r.nextReboot.Format(TimeFormat),
			)
		}
	}
}
