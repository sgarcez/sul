package uploader

import (
	"io"
	"log"

	strava "github.com/strava/go.strava"
)

// Uploader communicates with the Strava API
type Uploader struct {
	client  *strava.Client
	service *strava.UploadsService
}

// NewUploader initialises an Uploader instance
func NewUploader(token string) *Uploader {
	u := Uploader{}
	u.client = strava.NewClient(token)
	u.service = strava.NewUploadsService(u.client)
	return &u
}

// Upload creates and Activity from a file
func (u *Uploader) Upload(fname string, f io.Reader) *int64 {
	ft := strava.FileDataTypes.FIT
	resp, err := u.service.Create(ft, fname, f).Private().Do()
	if err != nil {
		if e, ok := err.(strava.Error); ok && e.Message == "Authorization Error" {
			log.Printf("%s - Auth error. Make sure your token has 'write' permissions.", fname)
		} else {
			log.Printf("%s - %s", fname, err)
		}
		return nil
	}

	uploadSummary, err := u.service.Get(resp.Id).Do()
	if err != nil {
		// TODO: parse error and sleep if activity isn't ready yet
		log.Printf("%s - %s", fname, err)
		return nil
	}

	log.Printf(
		"%s - Activity created, you can view it at http://www.strava.com/activities/%d",
		fname, uploadSummary.ActivityId)

	return &uploadSummary.ActivityId
}
