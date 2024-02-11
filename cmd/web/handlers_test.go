package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"


	"snippetbox.msp.net/internal/assert"
)


func TestPing(t *testing.T) {
	
	// Initlaize a new HTTPTEST Response recorder
	rr := httptest.NewRecorder()

	// Init a dummy 
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Call the route we want to test.
	ping(rr, r)
	
	// Grab the result from the ping 
	rs := rr.Result() 

	// see if we get a 200 code back
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	// Read the body of the response 
	body,err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
