package fbtron

import (
  "bytes"
  "net/http"
  "net/http/httptest"
  "testing"
  "time"
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
  recvchan := make(chan UserAction, 3)
  sendchan := make(chan Simulation, 3)

  server := httptest.NewServer(http.HandlerFunc(AddPlayers))
  defer server.Close()
  inchan = sendchan  // For the server
  outchan = recvchan

  // The server will expect these. This is a little messy.
  sendchan <- Simulation{}
  sendchan <- Simulation{}
  sendchan <- Simulation{}
  resp, err := http.Post(server.URL, "application/x-www-form-urlencoded",
                         bytes.NewBufferString("1234=1&2345=2&3456=3"))
  if err != nil {
    t.Errorf("AddPlayers: error %s", err)
  }

  // Make sure recvchan gets the right message
  expected_msgs := 3
  for expected_msgs > 0 {
    select {
    case msg := <-recvchan:
      if msg.action != ACTION_ADD {
        t.Errorf("AddPlayers: expected message action ACTION_ADD, got: %d",
                 msg.action)
        return
      }
      expected_msgs--
    case <-time.After(5 * time.Second):
      t.Errorf("AddPlayers: expected 3 messages, still waiting for %d more",
               expected_msgs)
    }
  }

  // Make sure we got an OK
  if resp.StatusCode != 200 {
    t.Errorf("AddPlayers: expected response code 200, got %d", resp.StatusCode)
  }
}

func TestRemovePlayers(t *testing.T) {
  recvchan := make(chan UserAction, 3)
  sendchan := make(chan Simulation, 3)

  server := httptest.NewServer(http.HandlerFunc(RemovePlayers))
  defer server.Close()
  inchan = sendchan  // For the server
  outchan = recvchan

  // The server will expect these. This is a little messy.
  sendchan <- Simulation{}
  sendchan <- Simulation{}
  sendchan <- Simulation{}
  resp, err := http.Post(server.URL, "application/x-www-form-urlencoded",
                         bytes.NewBufferString("1234=1&2345=2&3456=3"))
  if err != nil {
    t.Errorf("RemovePlayers: error %s", err)
  }

  // Make sure recvchan gets the right message
  expected_msgs := 3
  for expected_msgs > 0 {
    select {
    case msg := <-recvchan:
      if msg.action != ACTION_REM {
        t.Errorf("RemovePlayers: expected message action ACTION_REM, got: %d",
                 msg.action)
        return
      }
      expected_msgs--
    case <-time.After(5 * time.Second):
      t.Errorf("RemovePlayers: expected 3 messages, still waiting for %d more",
               expected_msgs)
    }
  }

  // Make sure we got an OK
  if resp.StatusCode != 200 {
    t.Errorf("RemovePlayers: expected response code 200, got %d",
             resp.StatusCode)
  }
}
