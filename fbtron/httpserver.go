package fbtron

import (
  "encoding/json"
  "flag"
  "fmt"
  "html/template"
  "net/http"
  "strconv"
)

const (
  ACTION_PING = iota
  ACTION_ADD  = iota
  ACTION_REM  = iota
  ACTION_QUIT = iota
)

type JsonData struct {
  Num_seasons   int
  Top_players   PlayerSlice
  Teams         []Team
}

type UserAction struct {
  action      int
  player_id   string
  team_id     int
}

var http_port = flag.Int("http_port", 8888, "Port to start http server on")

var inchan <-chan Simulation
var outchan chan<- UserAction

// RunWebServer starts a UI interface thread.
func RunWebServer(in <-chan Simulation, out chan<- UserAction) {
  inchan = in
  outchan = out

  http.HandleFunc("/", MainPage)
  http.HandleFunc("/data", GetData)
  http.HandleFunc("/add", AddPlayers)
  http.HandleFunc("/rem", RemovePlayers)
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
  outchan <- UserAction{ACTION_PING, "", 0}
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
  SendPlayerActions(w, req, ACTION_ADD)
}

func RemovePlayers(w http.ResponseWriter, req *http.Request) {
  SendPlayerActions(w, req, ACTION_REM)
}

func SendPlayerActions(w http.ResponseWriter, req *http.Request, action int) {
  req.ParseForm()

  // Send a message for every player
  msgs_sent := 0
  for player_id, team_ids := range req.PostForm {
    if len(team_ids) != 1 {
      fmt.Println("not enough team_ids for player_id ", player_id)
      continue
    }
    team_id, err := strconv.Atoi(team_ids[0])
    if err != nil {
      fmt.Println("error converting team_id ", string(team_ids[0]), ": ", err)
      continue
    }

    msgs_sent++
    outchan <- UserAction{action, string(player_id), team_id}
    _ = <-inchan  // TODO: Don't throw these away
  }

  if msgs_sent == 0 {
    http.Error(w, "No players to add/del", 500)
  }

  // TODO: Expect response?
  fmt.Printf("Got %d add/del messages\n", msgs_sent)
}
