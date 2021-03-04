/*main
  cyclesooon
  主程式進入點
*/
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/weiting4github/go-api-sooon/docs"

	_ "github.com/joho/godotenv/autoload"
)

// func TestPlayground(t *testing.T) {
// 	// The setupServer method, that we previously refactored
// 	// is injected into a test server
// 	ts := httptest.NewServer(setupServer())
// 	// Shut down the server and block until all requests have gone through
// 	defer ts.Close()

// 	// Make a request to our server with the {base url}/ping ts.URL like http://127.0.0.1:57040
// 	resp, err := http.Get(fmt.Sprintf("%s/dev/playground", ts.URL))

// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	if resp.StatusCode != 200 {
// 		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
// 	}

// 	val, ok := resp.Header["Content-Type"]

// 	// Assert that the "content-type" header is actually set
// 	if !ok {
// 		t.Fatalf("Expected Content-Type header to be set")
// 	}

// 	// Assert that it was set as expected
// 	if val[0] != "application/json; charset=utf-8" {
// 		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
// 	}
// }

func TestInit(t *testing.T) {
	// The setupServer method, that we previously refactored
	// is injected into a test server
	ts := httptest.NewServer(setupServer())
	// Shut down the server and block until all requests have gone through
	defer ts.Close()

	// Make a request to our server with the {base url}/ping ts.URL like http://127.0.0.1:57040
	resp, err := http.Post(fmt.Sprintf("%s/dev/init", ts.URL), "application/x-www-form-urlencoded", strings.NewReader("hash=5e72ee6bed764aed476bd7f369170b674a41518734df4b58975cef40fa09a224"))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	val, ok := resp.Header["Content-Type"]

	// Assert that the "content-type" header is actually set
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}

	// Assert that it was set as expected
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}
