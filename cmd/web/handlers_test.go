package main

import (
	"net/http"
	"testing"

	"snippetbox.msp.net/internal/assert"
)

func TestPing(t *testing.T) {
	
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	// Initlaize a new HTTPTEST Response recorder
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")

}
