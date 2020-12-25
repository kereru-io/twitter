package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/dghubble/oauth1"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type TwitterClient struct {
	Client *http.Client
}

const TwitterEndpoint string = "https://api.twitter.com/1.1/"
const TwitterUploadEndpoint string = "https://upload.twitter.com/1.1/"
var TwitterHTTPLogging bool = false

func NewTwitterClient(OauthConsumerKey, OauthConsumerSecret, OauthToken, OauthTokenSecret string) *TwitterClient {
	baseClient := &http.Client{
		Timeout: time.Second * 120,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			Proxy:        http.ProxyFromEnvironment,
		},
	}
	config := oauth1.NewConfig(OauthConsumerKey, OauthConsumerSecret)
	token := oauth1.NewToken(OauthToken, OauthTokenSecret)
	ctx := context.WithValue(oauth1.NoContext, oauth1.HTTPClient, baseClient)

	// client will automatically authorize http.Request's
	client := &TwitterClient{
		Client: oauth1.NewClient(ctx, config, token),
	}
	return client
}

func (client *TwitterClient) EncodeMultiPartForm(form url.Values) *bytes.Buffer {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for key, val := range form {
		for _, v := range val {
			_ = writer.WriteField(key, v)

		}
	}
	return body
}

func PrintRequest(r *http.Request) {
	if !TwitterHTTPLogging {
		return
	}
	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", dump)
}

func PrintResponse(r *http.Response) {
	if !TwitterHTTPLogging {
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", dump)
}

func (client *TwitterClient) DoRequest(req *http.Request) (interface{}, error) {
	PrintRequest(req)

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	PrintResponse(resp)

	//target := make(map[string]interface{})
	var target interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	if resp.Body != nil {
		err = decoder.Decode(&target)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}
	}
	//result := target.(map[string]interface{})
	return target, nil
}

func WriteMultiPartFile(writer *multipart.Writer, media []byte, fileName string) {
	part, err := writer.CreateFormFile("media", fileName)
	if err != nil {
		panic(err)
	}
	_, err = part.Write(media)
	if err != nil {
		panic(err)
	}
}

func WriteMultiPartField(writer *multipart.Writer, key string, value string) {
	err := writer.WriteField(key, value)
	if err != nil {
		panic(err)
	}
}

func (client *TwitterClient) EncodeForm(form url.Values) *bytes.Buffer {
	return bytes.NewBufferString(form.Encode())
}

func (client *TwitterClient) EncodeJson(form map[string]interface{}) *bytes.Buffer {
	body, err := json.Marshal(form)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(body)
}

func (client *TwitterClient) EncodeQuery(query url.Values) string {
	return query.Encode()
}
