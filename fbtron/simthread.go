package fbtron

import (
  "fmt"
  "runtime"
)

type Simulation struct {
  Teams           []Team
  Num_seasons     int
  Avail_players   []*Player
}

// RunSimulation is run as a goroutine. It receives information from main()
// about events, and it replies with its current status.
func RunSimulation(inchan <-chan string, outchan chan<- Simulation) {
  var sim Simulation

  sim.InitPlayers()
  sim.InitTeams()

  for {
    select {
    case msg := <-inchan:
      // TODO: How do we tell the simulator that players were drafted through
      // the web UI?
      switch msg {
      case "quitquitquit":
        break
      default:
        outchan <- sim
      }
    default:
      // No message ready, run a season
      sim.RunSeason()

      // Yield, in case something else needs to run
      runtime.Gosched()
    }
  }
}

// InitPlayers loads a set of players from the CSV files in the data directory.
func (sim *Simulation) InitPlayers() {
  batters := BuildPlayersFromCsv("data/steamer_hitters_2014_update.csv")
  pitchers := BuildPlayersFromCsv("data/steamer_pitchers_2014_update.csv")

  sim.Avail_players = make([]*Player, len(batters) + len(pitchers))
  copy(sim.Avail_players, batters)
  for n := range pitchers {
    sim.Avail_players[len(batters) + n] = pitchers[n]
  }

  fmt.Printf("Loaded %d players (%d batters, %d pitchers)\n",
             len(sim.Avail_players), len(batters), len(pitchers))
}

func (sim *Simulation) InitTeams() {
  // TODO: Create N teams.
}

// Run season simulates a single simulated season.
func (sim *Simulation) RunSeason() {
  sim.Num_seasons++

  // TODO: Perform the draft
  // TODO: Compare all pairs of teams
  // TODO: Award wins
  // TODO: End the season
}
