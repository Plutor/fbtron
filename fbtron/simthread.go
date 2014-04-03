package fbtron

import (
  "flag"
  "fmt"
  "math/rand"
  "runtime"
)

type PlayerSlice []*Player

type Simulation struct {
  Teams           []Team
  Num_seasons     int
  Avail_players   map[string]PlayerSlice
  All_players     PlayerSlice
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
  sim.AddPlayersToPositionLists(batters)
  sim.AddPlayersToPositionLists(pitchers)

  //sim.All_players = make(PlayerSlice, len(batters) + len(pitchers))
  sim.All_players = append(batters, pitchers...)
}

func (sim *Simulation) AddPlayersToPositionLists(players PlayerSlice) {
  if sim.Avail_players == nil {
    sim.Avail_players = make(map[string]PlayerSlice)
  }

  num_players := 0
  num_player_pos := 0

  for n := range players {
    num_players++
    for _, pos := range players[n].positions {
      num_player_pos++
      if sim.Avail_players[pos] == nil {
        sim.Avail_players[pos] = make(PlayerSlice, 0)
      }

      sim.Avail_players[pos] = append(sim.Avail_players[pos], players[n])
    }
  }
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
      p := sim.RandomAvailablePlayer(pos)
      if p == nil {
        // None available! BIG PROBLEM!
        // TODO: What do we do?
        break
      }

      // Add to the team
      team.AddPlayer(p, false)
    }
  }
}

// RandomAvailablePlayerIndex returns the index of a random available player
// that plays the given position.
func (sim *Simulation) RandomAvailablePlayer(position string) *Player {
  allplayers := sim.Avail_players[position]
  if len(allplayers) == 0 {
    return nil
  }

  // TODO: Use weighted randomness, picking better players more often.
  p := allplayers[rand.Intn(len(allplayers))]

  // Remove from any position list
  for _, pos := range p.positions {
    for n := range sim.Avail_players[pos] {
      if sim.Avail_players[pos][n] == p {
        sim.Avail_players[pos] = append(sim.Avail_players[pos][:n],
                                        sim.Avail_players[pos][n+1:]...)
        break
      }
    }
  }

  return p
}

// ScoreSeason compares each team to each other team. For each stat, the team
// with the greater value is awarded a win (ties are ignored).
func (sim *Simulation) ScoreSeason() {
  for a := range sim.Teams {
    for b := 0; b < a; b++ {
      for stat := range stat_types {
        diff := sim.Teams[a].GetStat(stat) - sim.Teams[b].GetStat(stat)
        if diff > 0 {
          sim.Teams[a].wins++
        } else if diff < 0 {
          sim.Teams[b].wins++
        }
      }
    }
  }
}

// EndSeason releases all non-keeper players (which implicitly credits them with
// their team's wins) and adds them back to the available players pool.
func (sim *Simulation) EndSeason() {
  for n := range sim.Teams {
    released_players := sim.Teams[n].Release()
    sim.AddPlayersToPositionLists(released_players)

    sim.Teams[n].wins = 0
  }
}
