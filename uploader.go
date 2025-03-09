package tus

import (
	"bytes"
	"net/http"
)

type Uploader struct {
	client     *Client
	url        string
	upload     *Upload
	offset     int64
	aborted    bool
	uploadSubs []chan Upload
	notifyChan chan bool
	response   *http.Response
}

// NotifyUploadProgress subscribes to progress updates.
func (u *Uploader) NotifyUploadProgress(c chan Upload) {
	u.uploadSubs = append(u.uploadSubs, c)
}

// Abort aborts the upload process.
// It doesn't abort the current chunk, only the remaining.
func (u *Uploader) Abort() {
	u.aborted = true
}

// IsAborted returns true if the upload was aborted.
func (u *Uploader) IsAborted() bool {
	return u.aborted
}

// Url returns the upload url.
func (u *Uploader) Url() string {
	return u.url
}

// Response returns the pre response url.
func (u *Uploader) Response() *http.Response {
	return u.response
}

// Offset returns the current offset uploaded.
func (u *Uploader) Offset() int64 {
	return u.offset
}

// Upload uploads the entire body to the server.
func (u *Uploader) Upload() error {
	for u.offset < u.upload.size && !u.aborted {
		err := u.UploadChunk()

		if err != nil {
			return err
		}
	}

	return nil
}

// UploadChunk uploads a single chunck.
func (u *Uploader) UploadChunk() error {
	data := make([]byte, u.client.Config.ChunkSize)

	_, err := u.upload.stream.Seek(u.offset, 0)

	if err != nil {
		return err
	}

	size, err := u.upload.stream.Read(data)

	if err != nil {
		return err
	}

	body := bytes.NewBuffer(data[:size])

	newOffset, res, err := u.client.uploadChunk(u.url, body, int64(size), u.offset, u.upload.size)

	if err != nil {
		return err
	}

	u.response = res

	u.offset = newOffset

	u.upload.updateProgress(u.offset)

	u.notifyChan <- true

	return nil
}

// Waits for a signal to broadcast to all subscribers
func (u *Uploader) broadcastProgress() {
	for _ = range u.notifyChan {
		for _, c := range u.uploadSubs {
			c <- *u.upload
		}
	}
}

// NewUploader creates a new Uploader.
func NewUploader(client *Client, url string, upload *Upload, offset int64) *Uploader {
	notifyChan := make(chan bool)

	uploader := &Uploader{
		client,
		url,
		upload,
		offset,
		false,
		nil,
		notifyChan,
		nil,
	}

	go uploader.broadcastProgress()

	return uploader
}
