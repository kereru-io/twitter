package api

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (client *TwitterClient) StatusesUserTimeline(user_id string) (interface{}, error) {
	path := TwitterEndpoint + "statuses/user_timeline.json"

	params := url.Values{}
	params.Set("user_id", user_id)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		panic(err)
	}
	req.URL.RawQuery = client.EncodeQuery(params)
	return client.DoRequest(req)
}

func (client *TwitterClient) StatusesUpdate(status string) (interface{}, error) {
	path := TwitterEndpoint + "statuses/update.json"

	params := url.Values{}
	params.Set("status", status)

	body := client.EncodeForm(params)

	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return client.DoRequest(req)
}

func (client *TwitterClient) StatusesUpdateWithMedia(status, mediaIDs string) (interface{}, error) {
	path := TwitterEndpoint + "statuses/update.json"

	params := url.Values{}
	params.Set("status", status)
	params.Set("media_ids", mediaIDs)

	body := client.EncodeForm(params)

	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return client.DoRequest(req)
}

func (client *TwitterClient) MediaUploadInit(totalBytes, mediaType, mediaCategory string) (interface{}, error) {
	path := TwitterUploadEndpoint + "media/upload.json"

	params := url.Values{}
	params.Set("command", "INIT")
	params.Set("total_bytes", totalBytes)
	params.Set("media_type", mediaType)
	params.Set("media_category", mediaCategory)

	body := client.EncodeForm(params)

	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return client.DoRequest(req)
}

func (client *TwitterClient) MediaUploadAppend(media []byte, mediaID, fileName, segmentIndex string) (interface{}, error) {
	path := TwitterUploadEndpoint + "media/upload.json"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	WriteMultiPartFile(writer, media, fileName)
	WriteMultiPartField(writer, "command", "APPEND")
	WriteMultiPartField(writer, "media_id", mediaID)
	WriteMultiPartField(writer, "segment_index", segmentIndex)

	err := writer.Close()
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return client.DoRequest(req)
}

func (client *TwitterClient) MediaUploadFinalize(mediaID string) (interface{}, error) {
	path := TwitterUploadEndpoint + "media/upload.json"

	params := url.Values{}
	params.Set("command", "FINALIZE")
	params.Set("media_id", mediaID)

	body := client.EncodeForm(params)

	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return client.DoRequest(req)
}

func (client *TwitterClient) MediaUploadStatus(mediaID string) (interface{}, error) {
	path := TwitterUploadEndpoint + "media/upload.json"

	params := url.Values{}
	params.Set("command", "STATUS")
	params.Set("media_id", mediaID)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		panic(err)
	}
	req.URL.RawQuery = client.EncodeQuery(params)
	return client.DoRequest(req)
}

func (client *TwitterClient) MediaMetadataCreate(mediaID, alt_text string) (interface{}, error) {
	path := TwitterUploadEndpoint + "media/metadata/create.json"

	params := map[string]interface{}{
		"media_id": mediaID,
		"alt_text": map[string]interface{}{
			"text": alt_text,
		},
	}
	body := client.EncodeJson(params)

	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	return client.DoRequest(req)
}

func (client *TwitterClient) AccountVerifyCredentials() (interface{}, error) {
	path := TwitterEndpoint + "account/verify_credentials.json"

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return client.DoRequest(req)
}

func (client *TwitterClient) HelpConfiguration() (interface{}, error) {
	path := TwitterEndpoint + "help/configuration.json"

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return client.DoRequest(req)
}

func (client *TwitterClient) ApplicationRateLimitStatus() (interface{}, error) {
	path := TwitterEndpoint + "application/rate_limit_status.json"

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return client.DoRequest(req)
}
