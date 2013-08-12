package server

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"strings"
)

var (
	once       sync.Once
	serverAddr string
	server     *httptest.Server
)

func startServer() {
	LogLevel = 1
	setupServer("")
	server = httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
}

type SanitizeTest struct {
	in, out string
}

type BlacklistTest struct {
	in, out []string
	method string
}

func TestBlacklist(t *testing.T) {
	once.Do(startServer)

	tests := []*BlacklistTest{
		{[]string{"x"}, []string{"x"}, "POST"},
		{[]string{"x"}, []string{"x"}, "POST"},
		{[]string{"x"}, []string{"x"}, "PUT"},
		{[]string{"y"}, []string{"x","y"}, "PUT"},
		{[]string{"a"}, []string{"a", "x", "y"}, "PUT"},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c", "x", "y"}, "PUT"},
	}

	for i, x := range tests {
		blacklistHttp(t, i, x.in, x.out, x.method)
	}
}

func blacklistHttp(t *testing.T, index int, in, out []string, method string) {
	values := url.Values{}

	for _, s := range in {
		values.Add("blacklist", s)
	}

	params := strings.NewReader(values.Encode())
	req, _ := http.NewRequest(method, fmt.Sprintf("http://%s/api/1.0/blacklist/", serverAddr), params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	r, err := client.Do(req)

	if err != nil {
		t.Fatalf("error posting: %s", err)
		return
	}

	if method == "POST" {
		if r.StatusCode != 201 {
			t.Fatalf("expected status code 201, got %d", r.StatusCode)
		} 
	} else {
		if r.StatusCode != 200 {
			t.Fatalf("expected status code 200, got %d", r.StatusCode)
		}
	}

	r, err = http.Get(fmt.Sprintf("http://%s/api/1.0/blacklist/", serverAddr))

	if err != nil {
		t.Fatalf("error getting: %s", err)
		return
	}

	var res []string
	err = json.NewDecoder(r.Body).Decode(&res)

	if err != nil {
		t.Fatal(err)
	}

	//Logf("res: %s", res)

	if len(res) != len(out) {
		t.Fatalf("%d != %d", len(res), len(out))
	}

	for i := 0; i < len(res); i++ {
		if res[i] != out[i] {
			t.Fatalf("%s != %s", res[i], out[i])
		}
	}
}

func TestSanitize(t *testing.T) {
	once.Do(startServer)
	blacklistHttp(t, 0, []string{"xxxx"}, []string{"xxxx"}, "POST")
	tests := []*SanitizeTest{
		{"foo", "foo"},
		{"foo xxxx", "foo ****"},
		{"foo uck", "foo uck"},
		{"foo fxxxx", "foo f****"},
	}

	for i, x := range tests {
		sanitizeHttp(t, i, x.in, x.out)
	}
}

func sanitizeHttp(t *testing.T, index int, in, out string) {
	values := url.Values{
		"text": {in},
	}

	r, err := http.Get(fmt.Sprintf("http://%s/api/1.0/sanitize/?%s", serverAddr, values.Encode()))

	if err != nil {
		t.Fatalf("error posting: %s", err)
		return
	}

	res := new(Response)
	err = json.NewDecoder(r.Body).Decode(res)

	if err != nil {
		t.Fatal(err)
	}

	Logf("res: %s", res)

	if r.StatusCode != 200 {
		t.Fatalf("expected status code 200, got %d", r.StatusCode)
	}

	if res.Text != out {
		t.Fatalf("#%d: expected %s, got %s", index, out, res.Text)
	}
}

func BenchmarkServer(b *testing.B) {
	in := []string{"a", "b", "c", "d", "e", "f", "g", "h","k", "l", "m", "n", "o", "p", "q", "r", "s", "t"}
	values := url.Values{}

	for _, s := range in {
		values.Add("blacklist", s)
	}

	params := strings.NewReader(values.Encode())
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/api/1.0/blacklist/", serverAddr), params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	_, err := client.Do(req)

	if err != nil {
		b.Fatalf("error posting: %s", err)
		return
	}

	values = url.Values{
		"text": in,
	}
	uri := fmt.Sprintf("http://%s/api/1.0/sanitize/?%s", serverAddr, values.Encode())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := http.Get(uri)

		if err != nil {
			b.Fatalf("error posting: %s", err)
		}
	}
}
