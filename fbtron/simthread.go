package fbtron

import (
  "flag"
  "fmt"
  "math/rand"
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
  sim.InitTeams(POSITIONS)

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
func (sim *Simulation) InitTeams(positions map[string]int) {
  sim.Teams = make([]Team, *num_teams)
  for n := 0; n < *num_teams; n ++ {
    sim.Teams[n] = Team {
      name: fmt.Sprintf("Team %d", n),
    }
    sim.Teams[n].SetPositions(positions)
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
  for n := range sim.Teams {
    team := sim.Teams[n]
    for {
      pos := team.GetOpenPosition()
      if pos == "" {
        break
      }

      // Choose a random available player
      pindex := sim.RandomAvailablePlayerIndex(pos)
      if pindex < 0 {
        // None available! BIG PROBLEM!
        // TODO: What do we do?
        break
      }

      // Add to the team, remove from available
      team.AddPlayer(sim.Avail_players[pindex], false)
      sim.Avail_players = append(sim.Avail_players[:pindex],
                                 sim.Avail_players[pindex+1:]...)
    }
  }
}

// RandomAvailablePlayerIndex returns the index of a random available player
// that plays the given position.
func (sim *Simulation) RandomAvailablePlayerIndex(position string) int {
  allindexes := sim.AllAvailablePlayersIndexes(position)
  if len(allindexes) == 0 {
    return -1
  }
  // TODO: Use weighted randomness, picking better players more often.
  return allindexes[rand.Intn(len(allindexes))]
}

// RandomAvailablePlayerIndex returns the indexes of all of the available
// players that play the given position.
func (sim *Simulation) AllAvailablePlayersIndexes(position string) []int {
  allindexes := make([]int, 0, len(sim.Avail_players))
  for n := range sim.Avail_players {
    for pos := range sim.Avail_players[n].positions {
      if sim.Avail_players[n].positions[pos] == position {
        allindexes = append(allindexes, n)
        break
      }
    }
  }
  return allindexes
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
