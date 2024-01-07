package user_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	resp, err := http.Post("http://localhost:8080/user/sign-in", "application/json", bytes.NewBuffer([]byte(`{"email": "test@test"}`)))
	
	if err != nil {
		t.Error(err.Error())
	}

  defer resp.Body.Close()
    
  respBody, err := io.ReadAll(resp.Body)

  if err == nil {
    strBody := string(respBody)
		fmt.Println(strBody)
    assert.Contains(t, strBody, "access_token", "response body must contain access_token")
  }
}