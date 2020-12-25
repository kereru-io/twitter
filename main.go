package twitter

import (
	"encoding/json"
	"github.com/kereru-io/twitter/api"
	"log"
	"strconv"
	"time"
)

func uploadFileProcessingInfo(processingInfo map[string]interface{}) (string, int64) {
	var ok bool

	state := processingInfo["state"].(string)

	log.Printf("State: %v", state)

	switch state {
	case "pending":
		fallthrough
	case "in_progress":
		if _, ok = processingInfo["check_after_secs"].(json.Number); !ok {
			log.Printf("Invalid value for check_after_secs: %v", processingInfo["check_after_secs"])
			return state, 0
		}
		checkAfterSecs := processingInfo["check_after_secs"].(json.Number)
		log.Printf("checkAfterSecs: %v", checkAfterSecs)

		sleepTime, err := processingInfo["check_after_secs"].(json.Number).Int64()
		if err != nil {
			panic(err)
		}
		log.Printf("sleepTime: %v", sleepTime)
		return state, sleepTime
	case "succeeded":
		return state, 0
	case "failure":
		return state, 0
	}
	return state, 0

}

func uploadFileInit(client *api.TwitterClient, size string, mimetype string, category string) string {
	var err error
	var ok bool
	var result interface{}

	result, err = client.MediaUploadInit(size, mimetype, category)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", result)

	if _, ok = result.(map[string]interface{}); !ok {
		log.Printf("Invalid value for result: %v", result)
		return ""
	}

	if _, ok = result.(map[string]interface{})["media_id_string"].(string); !ok {
		log.Printf("Invalid value for result: %v", result)
		return ""
	}

	return result.(map[string]interface{})["media_id_string"].(string)
}

func uploadFileAppend(client *api.TwitterClient, data []byte, media_id string, filename string, segment_index string) {
	var err error
	var result interface{}

	result, err = client.MediaUploadAppend(data, media_id, filename, segment_index)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", result)

	// this should return nil
	//if _ , ok = result.(n); !ok {
	//  log.Printf("Invalid value: %v", result)
	//  return
	//}
}

func uploadFileFinalize(client *api.TwitterClient, media_id string) (string, int64) {
	var err error
	var ok bool
	var result interface{}

	result, err = client.MediaUploadFinalize(media_id)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", result)

	if _, ok = result.(map[string]interface{}); !ok {
		log.Printf("Invalid value for result: %v", result)
		return "", 0
	}

	if _, ok = result.(map[string]interface{})["processing_info"].(map[string]interface{}); !ok {
		log.Printf("Invalid value for result processing_info: %v", result.(map[string]interface{})["processing_info"])
		return "succeeded", 0
	}

	if _, ok = result.(map[string]interface{})["processing_info"].(map[string]interface{})["state"].(string); !ok {
		log.Printf("Invalid value for result processing_info state: %v", result.(map[string]interface{})["processing_info"].(map[string]interface{})["state"].(string))
		return "", 0
	}
	return uploadFileProcessingInfo(result.(map[string]interface{})["processing_info"].(map[string]interface{}))

}

func uploadFileStatus(client *api.TwitterClient, media_id string) (string, int64) {
	var err error
	var ok bool
	var result interface{}

	result, err = client.MediaUploadStatus(media_id)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", result)

	if _, ok = result.(map[string]interface{}); !ok {
		log.Printf("Invalid value for result: %v", result)
		return "", 0
	}

	if _, ok = result.(map[string]interface{})["processing_info"].(map[string]interface{}); !ok {
		log.Printf("Invalid value for result processing_info: %v", result.(map[string]interface{})["processing_info"])
		return "", 0
	}

	if _, ok = result.(map[string]interface{})["processing_info"].(map[string]interface{})["state"].(string); !ok {
		log.Printf("Invalid value for result processing_info state: %v", result.(map[string]interface{})["processing_info"].(map[string]interface{})["state"])
		return "", 0
	}

	return uploadFileProcessingInfo(result.(map[string]interface{})["processing_info"].(map[string]interface{}))
}

func UploadFile(client *api.TwitterClient, filename string) string {
	const segmentSize int64 = 1024 * 1024 * 4

	file, _ := OpenFile(filename)
	size, _ := FileSizeString(file)
	log.Printf("Size: %v", size)
	category := FileCategory(file)
	log.Printf("Category: %v", category)
	mimetype, _ := FileMimeType(file)
	log.Printf("Mime Type: %v", mimetype)
	segment_count, _ := FileSegmentCount(segmentSize, file)
	log.Printf("Segment Count: %v", segment_count)

	media_id := uploadFileInit(client, size, mimetype, category)
	if media_id == "" {
		return ""
	}

	for i := int64(0); i < segment_count; i++ {
		data, _ := FileContents(i, segmentSize, file)
		data_length := len(data)
		segment_index := strconv.FormatInt(i, 10)
		log.Printf("Segment id: %v Index: %v Data length: %v", i, segment_index, data_length)

		uploadFileAppend(client, data, media_id, filename, segment_index)
	}

	status, sleepTime := uploadFileFinalize(client, media_id)

	for {
		switch status {
		case "succeeded":
			return media_id
		case "failed":
			return ""
		case "pending":
			fallthrough
		case "in_progress":
			if sleepTime > 0 {
				time.Sleep(time.Duration(sleepTime) * time.Second)
				status, sleepTime = uploadFileStatus(client, media_id)
			}
		default:
			return ""
		}
	}
}
