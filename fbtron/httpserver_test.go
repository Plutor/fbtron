package fbtron

import (
  "bytes"
  "net/http"
  "net/http/httptest"
  "net/url"
  "testing"
)

func TestGetData(t *testing.T) {
  var sim Simulation

  recvchan := make(chan UserAction, 1)
  sendchan := make(chan Simulation, 1)

  server := httptest.NewServer(http.HandlerFunc(GetData))
  defer server.Close()
  inchan = sendchan  // For the server
  outchan = recvchan

  sendchan <- sim

  resp, err := http.Get(server.URL)
  if err != nil {
    t.Errorf("GetData: error %s", err)
  }

  // Make sure we got an OK
  if resp.StatusCode != 200 {
    t.Errorf("GetData: expected response code 200, got %d", resp.StatusCode)
  }

  // TODO: Test actual content
}

func TestAddPlayers(t *testing.T) {
  recvchan := make(chan UserAction, 1)
  sendchan := make(chan Simulation, 1)

  server := httptest.NewServer(http.HandlerFunc(AddPlayers))
  defer server.Close()
  inchan = sendchan  // For the server
  outchan = recvchan

  expectedval := "{1234:1,2345:2,3456:3}"

  resp, err := http.Post(server.URL, "text/json",
                         bytes.NewBufferString(expectedval))
  if err != nil {
    t.Errorf("AddPlayers: error %s", err)
  }

  // Make sure we got an OK
  if resp.StatusCode != 200 {
    t.Errorf("AddPlayers: expected response code 200, got %d", resp.StatusCode)
  }

  // Make sure recvchan gets the right message
  // select {
  // case msg := <-recvchan:
  //   if msg != expectedval {
  //     t.Errorf("AddPlayers: expected message:\n%s\ngot:\n%s", expectedval, msg)
  //   }
  // default:
  //   t.Errorf("AddPlayers: expected message, got none")
  // }
}

func TestRemovePlayer(t *testing.T) {
  recvchan := make(chan UserAction, 1)
  sendchan := make(chan Simulation, 1)

  server := httptest.NewServer(http.HandlerFunc(RemovePlayer))
  defer server.Close()
  inchan = sendchan  // For the server
  outchan = recvchan

  var postvalues url.Values
  resp, err := http.PostForm(server.URL, postvalues)
  if err != nil {
    t.Errorf("RemovePlayer: error %s", err)
  }

  // Make sure we got an OK
  if resp.StatusCode != 200 {
    t.Errorf("RemovePlayer: expected response code 200, got %d",
             resp.StatusCode)
  }

  // TODO: Make sure recvchan gets the right request
}
