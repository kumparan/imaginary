package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

const fixtureFile = "testdata/large.jpg"

func TestSourceBodyMatch(t *testing.T) {
	u, _ := url.Parse("http://foo")
	req := &http.Request{Method: http.MethodPost, URL: u}
	source := NewBodyImageSource(&SourceConfig{})

	if !source.Matches(req) {
		t.Error("Cannot match the request")
	}
}

func TestBodyImageSource(t *testing.T) {
	var body []byte
	var err error

	source := NewBodyImageSource(&SourceConfig{})
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		body, err = source.GetImage(r)
		if err != nil {
			t.Fatalf("Error while reading the body: %s", err)
		}
		w.Write(body)
	}

	file, _ := os.Open(fixtureFile)
	r, _ := http.NewRequest(http.MethodPost, "http://foo/bar", file)
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	buf, _ := ioutil.ReadFile(fixtureFile)
	if len(body) != len(buf) {
		t.Error("Invalid response body")
	}
}

func TestBodyImageSource_JSON_base64(t *testing.T) {
	var body []byte
	var err error

	source := NewBodyImageSource(&SourceConfig{})
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		body, err = source.GetImage(r)
		if err != nil {
			t.Fatalf("Error while reading the body: %s", err)
		}
		w.Write(body)
	}

	supportedJSONField := struct {
		Base64 string `json:"base64"`
	}{
		Base64: "/9j/4AAQSkZJRgABAQEAYABgAAD//gA+Q1JFQVRPUjogZ2QtanBlZyB2MS4wICh1c2luZyBJSkcgSlBFRyB2ODApLCBkZWZhdWx0IHF1YWxpdHkK/9sAQwAIBgYHBgUIBwcHCQkICgwUDQwLCwwZEhMPFB0aHx4dGhwcICQuJyAiLCMcHCg3KSwwMTQ0NB8nOT04MjwuMzQy/9sAQwEJCQkMCwwYDQ0YMiEcITIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIy/8AAEQgAMgAyAwEiAAIRAQMRAf/EAB8AAAEFAQEBAQEBAAAAAAAAAAABAgMEBQYHCAkKC//EALUQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+v/EAB8BAAMBAQEBAQEBAQEAAAAAAAABAgMEBQYHCAkKC//EALURAAIBAgQEAwQHBQQEAAECdwABAgMRBAUhMQYSQVEHYXETIjKBCBRCkaGxwQkjM1LwFWJy0QoWJDThJfEXGBkaJicoKSo1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoKDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uLj5OXm5+jp6vLz9PX29/j5+v/aAAwDAQACEQMRAD8A+f6KKKACiiigAooooAKKKKACpIreaff5UbPsXc20ZwKjq3YXos3bfGZY2xujyAGx2OQePyoAhltp4FRpoZI1cZUspGf85H5ikit5p9/lRs+xdzbRnAq9NqoknSUW+CJpJ2VnyGZwMjp935enoSM02HUxHdyzPDujlwWhDAIcduQfl7AdhxmgCpLbTwKjTQyRq4ypZSM/5yPzFRVdv78Xg4iKEyvM2Wzln25xwMD5feqVABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFAH//2Q==",
	}
	bodyByte, err := json.Marshal(supportedJSONField)
	if err != nil {
		t.Fatalf("Error marshalling body: %s", err)
	}

	r, _ := http.NewRequest(http.MethodPost, "http://foo/bar", strings.NewReader(string(bodyByte)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	buf, _ := base64.StdEncoding.DecodeString(supportedJSONField.Base64)
	if len(body) != len(buf) {
		t.Error("Invalid response body")
	}
}

func TestBodyImageSource_JSON_base64WithDataURLs(t *testing.T) {
	var body []byte
	var err error

	source := NewBodyImageSource(&SourceConfig{})
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		body, err = source.GetImage(r)
		if err != nil {
			t.Fatalf("Error while reading the body: %s", err)
		}
		w.Write(body)
	}

	supportedJSONField := struct {
		Base64 string `json:"base64"`
	}{
		Base64: "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD//gA+Q1JFQVRPUjogZ2QtanBlZyB2MS4wICh1c2luZyBJSkcgSlBFRyB2ODApLCBkZWZhdWx0IHF1YWxpdHkK/9sAQwAIBgYHBgUIBwcHCQkICgwUDQwLCwwZEhMPFB0aHx4dGhwcICQuJyAiLCMcHCg3KSwwMTQ0NB8nOT04MjwuMzQy/9sAQwEJCQkMCwwYDQ0YMiEcITIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIy/8AAEQgAMgAyAwEiAAIRAQMRAf/EAB8AAAEFAQEBAQEBAAAAAAAAAAABAgMEBQYHCAkKC//EALUQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+v/EAB8BAAMBAQEBAQEBAQEAAAAAAAABAgMEBQYHCAkKC//EALURAAIBAgQEAwQHBQQEAAECdwABAgMRBAUhMQYSQVEHYXETIjKBCBRCkaGxwQkjM1LwFWJy0QoWJDThJfEXGBkaJicoKSo1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoKDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uLj5OXm5+jp6vLz9PX29/j5+v/aAAwDAQACEQMRAD8A+f6KKKACiiigAooooAKKKKACpIreaff5UbPsXc20ZwKjq3YXos3bfGZY2xujyAGx2OQePyoAhltp4FRpoZI1cZUspGf85H5ikit5p9/lRs+xdzbRnAq9NqoknSUW+CJpJ2VnyGZwMjp935enoSM02HUxHdyzPDujlwWhDAIcduQfl7AdhxmgCpLbTwKjTQyRq4ypZSM/5yPzFRVdv78Xg4iKEyvM2Wzln25xwMD5feqVABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFABRRRQAUUUUAFFFFAH//2Q==",
	}
	bodyByte, err := json.Marshal(supportedJSONField)
	if err != nil {
		t.Fatalf("Error marshalling body: %s", err)
	}

	r, _ := http.NewRequest(http.MethodPost, "http://foo/bar", strings.NewReader(string(bodyByte)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	buf, _ := base64.StdEncoding.DecodeString(strings.Split(supportedJSONField.Base64, ",")[1])
	if len(body) != len(buf) {
		t.Error("Invalid response body")
	}
}

func testReadBody(t *testing.T) {
	var body []byte
	var err error

	source := NewBodyImageSource(&SourceConfig{})
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		if !source.Matches(r) {
			t.Fatal("Cannot match the request")
		}

		body, err = source.GetImage(r)
		if err != nil {
			t.Fatalf("Error while reading the body: %s", err)
		}
		w.Write(body)
	}

	file, _ := os.Open(fixtureFile)
	r, _ := http.NewRequest(http.MethodPost, "http://foo/bar", file)
	w := httptest.NewRecorder()
	fakeHandler(w, r)

	buf, _ := ioutil.ReadFile(fixtureFile)
	if len(body) != len(buf) {
		t.Error("Invalid response body")
	}
}
