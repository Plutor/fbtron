package fbtron

import (
  "flag"
  "fmt"
  "runtime"
)

type Simulation struct {
  Teams           []Team
  Num_seasons     int
  Avail_players   []*Player
}

var POSITIONS map[string]int = map[string]int {
  "C":    1,
  "1B":   1,
  "2B":   1,
  "SS":   1,
  "3B":   1,
  "OF":   4,
  "Util": 2,
  "SP":   2,
  "RP":   2,
  "P":    4,
}

var num_teams = flag.Int("teams", 10, "Number of teams")

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
  batters := BuildPlayersFromCsv("data/steamer_hitters_2014_update.csv", "Util")
  pitchers := BuildPlayersFromCsv("data/steamer_pitchers_2014_update.csv", "P")

  sim.Avail_players = make([]*Player, len(batters) + len(pitchers))
  copy(sim.Avail_players, batters)
  for n := range pitchers {
    sim.Avail_players[len(batters) + n] = pitchers[n]
  }

  fmt.Printf("Loaded %d players (%d batters, %d pitchers)\n",
             len(sim.Avail_players), len(batters), len(pitchers))
}

// InitTeams creates a set of teams with empty rosters.
func (sim *Simulation) InitTeams() {
  sim.Teams = make([]Team, *num_teams)
  for n := 0; n < *num_teams; n ++ {
    sim.Teams[n] = Team {
      name: fmt.Sprintf("Team %d", n),
    }
    sim.Teams[n].SetPositions(POSITIONS)
  }
}

// Run season simulates a single simulated season.
func (sim *Simulation) RunSeason() {
  sim.Num_seasons++

  sim.DoDraft()
  sim.ScoreSeason()
  sim.EndSeason()
}

func (sim *Simulation) DoDraft() {
  // TODO: Perform the draft
}

func (sim *Simulation) ScoreSeason() {
  // TODO: Compare all pairs of teams and award wins
}

// EndSeason releases all non-keeper players (which implicitly credits them with
// their team's wins) and adds them back to the available players pool.
func (sim *Simulation) EndSeason() {
  for n := range sim.Teams {
    released_players := sim.Teams[n].Release()
    sim.Avail_players = append(sim.Avail_players, released_players...)
  }
}
