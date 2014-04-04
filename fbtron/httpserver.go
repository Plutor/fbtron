package fbtron

import (
  "encoding/json"
  "flag"
  "fmt"
  "html/template"
  "net/http"
)

var inchan <-chan Simulation
var outchan chan<- string
var http_port = flag.Int("http_port", 8888, "Port to start http server on")

// RunWebServer starts a UI interface thread.
func RunWebServer(in <-chan Simulation, out chan<- string) {
  inchan = in
  outchan = out

  http.HandleFunc("/", MainPage)
  http.HandleFunc("/data", GetData)

  fmt.Printf("Listening on http://localhost:%d\n", *http_port)
  err := http.ListenAndServe(fmt.Sprintf(":%d", *http_port), nil)
  if err != nil {
    fmt.Println("ListenAndServe: ", err)
  }
}

// MainPage builds a static page that includes some fun JavaScript.
func MainPage(w http.ResponseWriter, req *http.Request) {
  t, _ := template.ParseFiles("templates/home.tmpl")
  t.Execute(w, "")
}

// GetData builds JSON that represents the current state of the simulation. The
// JavaScript on the main page queries this data and displays it.
func GetData(w http.ResponseWriter, req *http.Request) {
  outchan <- "ping"
  sim_totals := <-inchan

  var bestplayer *Player
  var bestwpd float64
  for _, player := range sim_totals.All_players {
    if bestplayer == nil || player.WinsPerDraft() > bestwpd {
      bestplayer = player
      bestwpd = player.WinsPerDraft()
    }
  }

  // TODO: Make some sort of flatter struct for encoding.
  enc := json.NewEncoder(w)
  if err := enc.Encode(&sim_totals); err != nil {
    fmt.Println(err)
  }
}
