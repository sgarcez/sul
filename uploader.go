package main

import (
	"io"
	"log"

	strava "github.com/strava/go.strava"
)

// uploader communicates with the Strava API
type uploader struct {
	client  *strava.Client
	service *strava.UploadsService
}

// newUploader initialises an Uploader instance
func newUploader(token string) *uploader {
	u := uploader{}
	u.client = strava.NewClient(token)
	u.service = strava.NewUploadsService(u.client)
	return &u
}

// Upload creates a Strava Activity from a file
func (u *uploader) Upload(fname string, f io.Reader) (*int64, error) {
	ft := strava.FileDataTypes.FIT
	resp, err := u.service.Create(ft, fname, f).Private().Do()
	if err != nil {
		if e, ok := err.(strava.Error); ok && e.Message == "Authorization Error" {
			log.Printf("%s - Auth error. Make sure your token has 'write' permissions.", fname)
		} else {
			log.Printf("%s - %s", fname, err)
		}
		return nil, err
	}

	uploadSummary, err := u.service.Get(resp.Id).Do()
	if err != nil {
		// TODO: parse error and sleep if activity isn't ready yet
		log.Printf("%s - %s", fname, err)
		return nil, err
	}

	return &uploadSummary.ActivityId, nil
}
