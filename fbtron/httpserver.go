package fbtron

import (
  "encoding/json"
  "flag"
  "fmt"
  "html/template"
  "io/ioutil"
  "net/http"
)

type JsonData struct {
  Num_seasons   int
  Top_players   PlayerSlice
  Teams         []Team
}

var http_port = flag.Int("http_port", 8888, "Port to start http server on")

var inchan <-chan Simulation
var outchan chan<- string

// RunWebServer starts a UI interface thread.
func RunWebServer(in <-chan Simulation, out chan<- string) {
  inchan = in
  outchan = out

  http.HandleFunc("/", MainPage)
  http.HandleFunc("/data", GetData)
  http.HandleFunc("/add", AddPlayers)
  http.Handle("/static/",
      http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

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

  enc := json.NewEncoder(w)
  if err := enc.Encode(&JsonData {
                           Num_seasons: sim_totals.Num_seasons,
                           Top_players: sim_totals.TopPlayers(100),
                           Teams:       sim_totals.Teams,
                       }); err != nil {
    fmt.Println(err)
  }
}

func AddPlayers(w http.ResponseWriter, req *http.Request) {
  defer req.Body.Close()
  reqbody, err := ioutil.ReadAll(req.Body)

  if err != nil {
    http.Error(w, err.Error(), 500)
  }

  // TODO: Make sure the string is parseable JSON

  fmt.Printf("Got add post: %s\n", string(reqbody))
  outchan <- string(reqbody)
  // TODO: Expect response?
}

func RemovePlayer(w http.ResponseWriter, req *http.Request) {
}
