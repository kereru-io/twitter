package twitter

import (
	"image/gif"
	"io"
	"net/http"
	"os"
	"strconv"
)

func OpenFile(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func FileMimeType(file *os.File) (string, error) {
	buffer := make([]byte, 512)
	n, err := file.ReadAt(buffer, 0)
	if err != nil && err != io.EOF {
		return "", err
	}

	return http.DetectContentType(buffer[:n]), nil
}

func FileSizeString(file *os.File) (string, error) {
	stat, err := file.Stat()
	if err != nil {
		return "0", err
	}

	return strconv.FormatInt(stat.Size(), 10), nil
}

func FileSize(file *os.File) (int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func FileContents(index, size int64, file *os.File) ([]byte, error) {
	buffer := make([]byte, size)
	offset := size * index
	bytes, err := file.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buffer[:bytes], nil
}

func FileSegmentCount(size int64, file *os.File) (int64, error) {
	totalSize, err := FileSize(file)
	if err != nil {
		return 0, err
	}
	return (totalSize / size) + 1, nil
}
func FileCategory(file *os.File) string {
	mimeType, _ := FileMimeType(file)
	switch mimeType {
	case "video/mp4":
		return "tweet_video"
	case "image/png":
		fallthrough
	case "image/jpeg":
		fallthrough
	case "image/webp":
		return "tweet_image"
	case "image/gif":
		foo, err := gif.DecodeAll(file)
		if err != nil {
			return "tweet_image"
		}
		if len(foo.Image) > 1 {
			return "tweet_gif"
		}
		return "tweet_image"
	default:
		return ""
	}
}
