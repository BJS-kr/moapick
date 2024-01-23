package test_utils

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"os"
)

type Tester struct {
	AccessToken string
}

func (t Tester) GET(path string) *http.Response {
	req, err := http.NewRequest("GET", path, nil)

	setHeaders(req, t.AccessToken)

	if err != nil {
		panic(err.Error())
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err.Error())
	}

	return res
}

func (t Tester) POST(path string, body string) *http.Response {
	req, err := http.NewRequest("POST", path, makeRawBody(body))

	setHeaders(req, t.AccessToken)

	if err != nil {
		panic(err.Error())
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err.Error())
	}

	return res
}

func (t Tester) PUT(path string, body string) *http.Response {
	req, err := http.NewRequest("PUT", path, makeRawBody(body))

	setHeaders(req, t.AccessToken)

	if err != nil {
		panic(err.Error())
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err.Error())
	}

	return res
}

func (t Tester) PATCH(path string, rawBody string) *http.Response {
	req, err := http.NewRequest("PATCH", path, makeRawBody(rawBody))

	setHeaders(req, t.AccessToken)

	if err != nil {
		panic(err.Error())
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err.Error())
	}

	return res
}

func (t Tester) DELETE(path string) *http.Response {
	req, err := http.NewRequest("DELETE", path, nil)

	setHeaders(req, t.AccessToken)

	if err != nil {
		panic(err.Error())
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err.Error())
	}

	return res
}

func makeRawBody(rawBody string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(rawBody))
}

func setHeaders(req *http.Request, accessToken string) {
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
}

func GetRawDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", host, user, password, port, dbname))

	if err != nil {
		panic(err)
	}

	return db
}
