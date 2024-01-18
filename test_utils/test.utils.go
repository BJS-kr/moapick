package testutils

import (
	"io"
	"net/http"
)

type Tester struct {
	GET    func(path string) *http.Response
	PUT    func(path string, body io.Reader) *http.Response
	POST   func(path string, body io.Reader) *http.Response
	PATCH  func(path string, body io.Reader) *http.Response
	DELETE func(path string) *http.Response
}

func MakeHTTPTester(accessToken string) Tester {
	return Tester{
		GET: func(path string) *http.Response {
			req, err := http.NewRequest("GET", path, nil)

			setHeaders(req, accessToken)

			if err != nil {
				panic(err.Error())
			}

			res, err := http.DefaultClient.Do(req)

			if err != nil {
				panic(err.Error())
			}

			return res
		},
		POST: func(path string, body io.Reader) *http.Response {
			req, err := http.NewRequest("POST", path, body)

			setHeaders(req, accessToken)

			if err != nil {
				panic(err.Error())
			}

			res, err := http.DefaultClient.Do(req)

			if err != nil {
				panic(err.Error())
			}

			return res
		},
		PUT: func(path string, body io.Reader) *http.Response {
			req, err := http.NewRequest("POST", path, body)

			setHeaders(req, accessToken)

			if err != nil {
				panic(err.Error())
			}

			res, err := http.DefaultClient.Do(req)

			if err != nil {
				panic(err.Error())
			}

			return res
		},
		PATCH: func(path string, body io.Reader) *http.Response {
			req, err := http.NewRequest("POST", path, body)

			setHeaders(req, accessToken)

			if err != nil {
				panic(err.Error())
			}

			res, err := http.DefaultClient.Do(req)

			if err != nil {
				panic(err.Error())
			}

			return res
		},
		DELETE: func(path string) *http.Response {
			req, err := http.NewRequest("GET", path, nil)

			setHeaders(req, accessToken)

			if err != nil {
				panic(err.Error())
			}

			res, err := http.DefaultClient.Do(req)

			if err != nil {
				panic(err.Error())
			}

			return res
		},
	}
}

func setHeaders(req *http.Request, accessToken string) {
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
}
