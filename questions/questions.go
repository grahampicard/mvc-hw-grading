package questions

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type serverQuestion func(string, string) (bool, string, error)
type responseTester func(*http.Response) (bool, error)

func statusText(pass bool) string {
	if pass {
		return "✅ PASS"
	}
	return "❌ FAIL"
}

// TestAll ...
func TestAll(nickname string, rawURL string, showOutput bool) (int, int, error) {
	doLog := func(args ...interface{}) {
		if showOutput {
			fmt.Println(args...)
		}
	}
	numPass := 0
	numFail := 0
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return numPass, numFail, err
	}

	questions := []serverQuestion{
		indexIsUp,
		nicknameIsUp,
		nicknameMatchesExpected(nickname),
	}
	for _, question := range questions {
		passed, questionText, err2 := question(parsedURL.Scheme, parsedURL.Host)
		doLog(statusText(passed && (err2 == nil)), "-", questionText)
		if passed {
			numPass++
		} else {
			numFail++
		}
	}
	return numPass, numFail, err
}

func newClient() *http.Client {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	return netClient
}

func testStatusEquals(response *http.Response, err error, questionText string, expectedStatus int) (bool, string, error) {
	if err != nil {
		return false, questionText, err
	}
	if response.StatusCode == expectedStatus {
		return true, questionText, nil
	}
	return false, questionText, nil
}

func readResponseBody(response *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	return bodyString, err
}

func testBodyEquals(response *http.Response, err error, questionText string, expectedBody string) (bool, string, error) {
	if err != nil {
		return false, questionText, err
	}
	dump, err2 := readResponseBody(response)
	if err2 != nil {
		return false, questionText, err
	}
	body := strings.Trim(string(dump), " ")
	if body == expectedBody {
		return true, questionText, nil
	}
	return false, questionText, nil
}

func testResponse(response *http.Response, err error, questionText string, testFunc responseTester) (bool, string, error) {
	if err != nil {
		return false, questionText, nil
	}
	result, err := testFunc(response)
	if result && err == nil {
		return true, questionText, nil
	}
	return false, questionText, nil
}

func getAndCheckFunction(scheme string, host string, urlPath string, query url.Values, questionText string, testFunc responseTester) (bool, string, error) {
	parsedURL := url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     urlPath,
		RawQuery: query.Encode(),
	}
	netClient := newClient()
	response, err := netClient.Get(parsedURL.String())
	return testResponse(response, err, questionText, testFunc)
}

func getAndCheckStatus(scheme string, host string, urlPath string, query url.Values, questionText string, expectedStatus int) (bool, string, error) {
	parsedURL := url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     urlPath,
		RawQuery: query.Encode(),
	}
	netClient := newClient()
	response, err := netClient.Get(parsedURL.String())
	return testStatusEquals(response, err, questionText, expectedStatus)
}

func getAndCheckBody(scheme string, host string, urlPath string, query url.Values, questionText string, expectedBody string, exact bool) (bool, string, error) {
	testFunc := func(response *http.Response) (bool, error) {
		body, err := readResponseBody(response)
		if err != nil {
			return false, err
		}

		if (exact && body == expectedBody) || (!exact && strings.Contains(body, expectedBody)) {
			return true, nil
		}
		return false, nil
	}
	return getAndCheckFunction(
		scheme,
		host,
		urlPath,
		query,
		questionText,
		testFunc,
	)
}

func indexIsUp(scheme string, baseURL string) (bool, string, error) {
	return getAndCheckStatus(
		scheme,
		baseURL,
		"/",
		url.Values{},
		"Your return a 200 status code at /",
		http.StatusOK,
	)
}

func nicknameIsUp(scheme string, baseURL string) (bool, string, error) {
	return getAndCheckStatus(
		scheme,
		baseURL,
		"/nickname",
		url.Values{},
		"Your return a 200 status code at /nickname",
		http.StatusOK,
	)
}

func nicknameMatchesExpected(nickname string) serverQuestion {
	return func(scheme string, baseURL string) (bool, string, error) {
		return getAndCheckBody(
			scheme,
			baseURL,
			"/nickname",
			url.Values{},
			fmt.Sprintf("Your response at /nickname includes your nickname: %s", nickname),
			nickname,
			false,
		)
	}
}

func debugHTML(n *html.Node) {
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		log.Fatalf("Render error: %s", err)
	}
	fmt.Println(buf.String())
}
