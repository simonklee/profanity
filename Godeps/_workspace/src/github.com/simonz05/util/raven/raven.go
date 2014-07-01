package raven

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	UDP_TEMPLATE = "%s\n\n%s"
)

var (
	maxFollows int = 10
)

type SentryTransport interface {
	Send(packet []byte, timestamp time.Time) (response string, err error)
}

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type HttpSentryTransport struct {
	PublicKey string
	URL       *url.URL
	Project   string
	Client    HttpClient
}

type UdpSentryTransport struct {
	PublicKey string
	URL       *url.URL
	Client    net.Conn
}

func (self *UdpSentryTransport) Send(packet []byte, timestamp time.Time) (response string, err error) {
	authHeader := AuthHeader(timestamp, self.PublicKey)
	udp_msg := fmt.Sprintf(UDP_TEMPLATE, authHeader, string(packet))
	self.Client.Write([]byte(udp_msg))

	return "", nil
}

func (self *HttpSentryTransport) Send(packet []byte, timestamp time.Time) (response string, err error) {
	apiURL := self.URL
	apiURL.User = nil

	// Append slash to prevent 301 redirect
	location := strings.TrimRight(apiURL.String(), "/") + "/"

	// for loop to follow redirects
	followCounter := 0
	for {
		buf := bytes.NewBuffer(packet)
		req, err := http.NewRequest("POST", location, buf)
		if err != nil {
			return "", err
		}

		authHeader := AuthHeader(timestamp, self.PublicKey)
		req.Header.Add("X-Sentry-Auth", authHeader)
		req.Header.Add("Content-Type", "application/octet-stream")
		req.Header.Add("Connection", "close")
		req.Header.Add("Accept-Encoding", "identity")

		resp, err := self.Client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 301 {
			// set the location to the new one to retry on the next iteration
			location = resp.Header["Location"][0]
			followCounter++
			if followCounter >= maxFollows {
				return "", fmt.Errorf("Was redirected more than %d times, giving up", maxFollows)
			}
		} else {
			// We want to return an error for anything that's not a
			// straight HTTP 200
			if resp.StatusCode != 200 {
				body, _ := ioutil.ReadAll(resp.Body)
				return string(body), errors.New(resp.Status)
			}
			body, _ := ioutil.ReadAll(resp.Body)
			return string(body), nil
		}
	}
	// should never get here
	panic("send broke out of loop")
}

type Client struct {
	URL       *url.URL
	PublicKey string
	SecretKey string
	Project   string
	Logger    string

	sentryTransport SentryTransport
}

type sentryRequest struct {
	Project   string                 `json:"project"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Logger    string                 `json:"logger"`
	Extra     map[string]interface{} `json:"extra"`
}

type sentryResponse struct {
	ResultId string `json:"result_id"`
}

// Template for the X-Sentry-Auth header
const xSentryAuthTemplate = "Sentry sentry_version=2.0, sentry_client=raven-go/0.1, sentry_timestamp=%v, sentry_key=%v"

// An iso8601 timestamp without the timezone. This is the format Sentry expects.
const iso8601 = "2006-01-02T15:04:05"

// NewClient creates a new client for a server identified by the given dsn
// A dsn is a string in the form:
//	{PROTOCOL}://{PUBLIC_KEY}:{SECRET_KEY}@{HOST}/{PATH}{PROJECT_ID}
// eg:
//	http://abcd:efgh@sentry.example.com/sentry/project1
func NewClient(dsn string, logger string) (self *Client, err error) {
	var sentryTransport SentryTransport

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	basePath := path.Dir(u.Path)
	project := path.Base(u.Path)
	if u.User == nil {
		return nil, errors.New(fmt.Sprintf("raven: DSN parse ERR. Missing `user` part in DSN %s", dsn))
	}

	publicKey := u.User.Username()
	secretKey, _ := u.User.Password()
	u.Path = basePath

	check := func(req *http.Request, via []*http.Request) error {
		fmt.Printf("%+v", req)
		return nil
	}

	switch {
	case u.Scheme == "udp":
		udp_conn, udp_err := net.Dial("udp", u.Host)
		if udp_err != nil {
			return nil, udp_err
		}
		sentryTransport = &UdpSentryTransport{URL: u,
			Client:    udp_conn,
			PublicKey: publicKey}
	case u.Scheme == "https":
		httpClient := &http.Client{
			Transport:     nil,
			CheckRedirect: check,
			Jar:           nil}
		u.Path = path.Join(u.Path, "/api/"+project+"/store/")
		sentryTransport = &HttpSentryTransport{
			URL:       u,
			Client:    httpClient,
			Project:   project,
			PublicKey: publicKey}
	case u.Scheme == "http":
		httpClient := &http.Client{
			Transport:     nil,
			CheckRedirect: check,
			Jar:           nil}
		u.Path = path.Join(u.Path, "/api/"+project+"/store/")
		sentryTransport = &HttpSentryTransport{
			URL:       u,
			Client:    httpClient,
			Project:   project,
			PublicKey: publicKey}
	default:
		return nil, fmt.Errorf("Invalid protocol specified: %s", u.Scheme)
	}

	return &Client{URL: u, PublicKey: publicKey, SecretKey: secretKey,
		sentryTransport: sentryTransport, Project: project, Logger: logger}, nil
}

func (client Client) Error(v ...interface{}) error {
	return client.CaptureMessage(fmt.Sprint(v...), nil)
}

func (client Client) Errorf(format string, v ...interface{}) error {
	return client.CaptureMessage(fmt.Sprintf(format, v), nil)
}

func (client Client) Errorln(v ...interface{}) error {
	return client.CaptureMessage(fmt.Sprintln(v...), nil)
}

// CaptureMessage sends a message to the Sentry server.
func (self *Client) CaptureMessage(message string, extra map[string]interface{}) (err error) {
	timestamp := time.Now().UTC()
	timestampStr := timestamp.Format(iso8601)

	packet := sentryRequest{
		Project:   self.Project,
		Message:   message,
		Timestamp: timestampStr,
		Level:     "error",
		Logger:    self.Logger,
		Extra:     extra,
	}

	buf := new(bytes.Buffer)
	b64Encoder := base64.NewEncoder(base64.StdEncoding, buf)
	writer := zlib.NewWriter(b64Encoder)
	jsonEncoder := json.NewEncoder(writer)

	if err := jsonEncoder.Encode(packet); err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	err = b64Encoder.Close()
	if err != nil {
		return err
	}

	_, err = self.sentryTransport.Send(buf.Bytes(), timestamp)
	if err != nil {
		return err
	}

	return nil
}

/* Compute the Sentry authentication header */
func AuthHeader(timestamp time.Time, publicKey string) string {
	return fmt.Sprintf(xSentryAuthTemplate, timestamp.Unix(),
		publicKey)
}
