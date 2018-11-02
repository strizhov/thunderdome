// Copyright (C) 2018

package thunderdomeserver

import (
	"context"
	"testing"
)

var (
	ctx                   = context.Background()
	sameAddress           = ":54800"
	uploadServerAddress   = ":54812"
	downloadServerAddress = ":54813"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("Error is missing!!")
	}
}

func TestWrongServerAddress(t *testing.T) {
	// error for incorrect address
	_, _, err := NewServerTest(t, "blah", "blah")
	assertError(t, err)

	// error for same address
	_, server, err := NewServerTest(t, sameAddress, sameAddress)
	assertNoError(t, err)
	err = server.Run(ctx)
	assertError(t, err)
}

//TODO: implement actual network connection to the server

/////////////////////////////////////////////////////////////////////////////
// Helpers
/////////////////////////////////////////////////////////////////////////////

type ServerTest struct {
}

func NewServerTest(tb testing.TB, upload_server, download_server string) (
	*ServerTest, *Server, error) {
	t := &ServerTest{}

	srv, err := t.New(upload_server, download_server)
	return t, srv, err
}

func (t *ServerTest) New(upload_server, download_server string) (
	s *Server, err error) {
	return New(upload_server, download_server)
}
