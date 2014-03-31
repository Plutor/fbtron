package fbtron

import (
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

func (sim *Simulation) InitPlayers() {
  // TODO: Load player data from CSV files
}

// Run season simulates a single simulated season.
func (sim *Simulation) RunSeason() {
  sim.Num_seasons++

  // TODO: Perform the draft
  // TODO: Compare all pairs of teams
  // TODO: Award wins
  // TODO: End the season
}
