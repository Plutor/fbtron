package fbtron

import (
  "flag"
  "fmt"
  "math/rand"
  "runtime"
  "sort"
)

type PlayerSlice []*Player
type PlayerSet map[string]*Player

type Simulation struct {
  Teams           []Team
  Num_seasons     int
  Avail_players   map[string]PlayerSet
  All_players     PlayerSet
}

var POSITIONS map[string]int = map[string]int {
  "C":  1,
  "1B": 1,
  "2B": 1,
  "SS": 1,
  "3B": 1,
  "OF": 4,
  "B":  2,  // Generic batter, usually "Util" in Yahoo FBB.
  "SP": 2,
  "RP": 2,
  "P":  4,
}

var num_teams = flag.Int("teams", 10, "Number of teams")

// RunSimulation is run as a goroutine. It receives information from main()
// about events, and it replies with its current status.
func RunSimulation(inchan <-chan UserAction, outchan chan<- Simulation) {
  var sim Simulation

  sim.InitPlayers()
  sim.InitTeams(POSITIONS)

  for {
    select {
    case msg := <-inchan:
      switch msg.action {
      case ACTION_QUIT:
        break
      case ACTION_ADD:
        sim.AddKeeper(msg.player_id, msg.team_id)
      case ACTION_REM:
        sim.RemoveKeeper(msg.player_id, msg.team_id)
      }

      // Any message gets the response of the current status
      outchan <- sim
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
  batters := BuildPlayersFromCsv("data/steamer_hitters_2014_update.csv", "B")
  pitchers := BuildPlayersFromCsv("data/steamer_pitchers_2014_update.csv", "P")
  sim.AddPlayersToPositionLists(batters)
  sim.AddPlayersToPositionLists(pitchers)

  sim.All_players = make(PlayerSet)
  for _, p := range append(batters, pitchers...) {
    sim.All_players[p.ID] = p
  }
}

func (sim *Simulation) AddPlayersToPositionLists(players []*Player) {
  if sim.Avail_players == nil {
    sim.Avail_players = make(map[string]PlayerSet, len(players))
  }

  num_players := 0
  num_player_pos := 0

  for _, player := range players {
    num_players++
    for _, pos := range player.Positions {
      num_player_pos++
      if sim.Avail_players[pos] == nil {
        sim.Avail_players[pos] = make(PlayerSet, 0)
      }

      sim.Avail_players[pos][player.ID] = player
    }
  }
}

// InitTeams creates a set of teams with empty rosters.
func (sim *Simulation) InitTeams(positions map[string]int) {
  sim.Teams = make([]Team, *num_teams)
  for n := 0; n < *num_teams; n ++ {
    sim.Teams[n] = Team {
      Name: fmt.Sprintf("Team %d", n),
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
  var player *Player
  pindex := rand.Intn(len(allplayers))
  for _, p := range allplayers {
    if pindex == 0 {
      player = p
      break
    }
    pindex--
  }

  // Remove from any position list
  for _, pos := range player.Positions {
    delete(sim.Avail_players[pos], player.ID)
  }

  return player
}

// ScoreSeason compares each team to each other team. For each stat, the team
// with the greater value is awarded a win (ties are ignored).
func (sim *Simulation) ScoreSeason() {
  stat_cache := make(map[string]map[int]float64, len(stat_types))
  for stat := range stat_types {
    stat_cache[stat] = make(map[int]float64, len(sim.Teams))
    for n, team := range sim.Teams {
      stat_cache[stat][n] = team.GetStat(stat)
    }
  }

  var astat, diff float64
  for _, cache := range stat_cache {
    for a := range sim.Teams {
      astat = cache[a]
      for b := 0; b < a; b++ {
        diff = astat - cache[b]
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

// Merges the Num_seasons and All_players from the passed simulation with this
// one. This is used for summing all of the simulation threads for reporting
// purposes.
func (sim *Simulation) Merge(other *Simulation) {
  sim.Num_seasons += other.Num_seasons

  // Merge the player stats
  if sim.All_players == nil {
    sim.All_players = make(PlayerSet)
  }
  for pid, player := range other.All_players {
    // TODO - This won't work right for randomized ids. We aren't guaranteed
    // they will be the same across threads.
    if sim.All_players[pid] == nil {
      // Create
      player_copy := *player
      sim.All_players[pid] = &player_copy
    } else {
      sim.All_players[pid].Num_seasons += player.Num_seasons
      sim.All_players[pid].Total_wins += player.Total_wins
    }
  }

  // Merge the team stats
  if sim.Teams == nil {
    sim.Teams = make([]Team, *num_teams)
    for n, team := range other.Teams {
      // Create
      sim.Teams[n] = team
    }
  } else {
    // TODO: Merge the stats (once a team has stats to merge)
  }
}

// Len(), Less(), and Swap() make PlayerSlice sort()able by WinsPerDraft()
func (ps PlayerSlice) Len() int {
  return len(ps)
}
func (ps PlayerSlice) Less(i, j int) bool {
  return ps[i].WinsPerDraft() > ps[j].WinsPerDraft()
}
func (ps PlayerSlice) Swap(i, j int) {
  ps[i], ps[j] = ps[j], ps[i]
}

func (sim *Simulation) TopPlayers(num int) PlayerSlice {
  rv := make(PlayerSlice, 0, len(sim.All_players))
  for _, player := range sim.All_players {
    // Only include players who are actually available.
    for _, team := range sim.Teams {
      if team.HasPlayer(player.ID, true) {
        goto NextPlayer
      }
    }
    rv = append(rv, player)
    NextPlayer:
  }
  sort.Sort(rv)

  if num <= 0 || num >= len(rv) {
    return rv
  }
  return rv[:num]
}

// AddKeeper adds a keeper to the specified team
func (sim *Simulation) AddKeeper(player_id string, team_id int) {
  if team_id >= len(sim.Teams) {
    fmt.Println("Couldn't find team")
    return
  }
  team := sim.Teams[team_id]

  // Find the player
  var player *Player
  for _, players := range sim.Avail_players {
    for _, p := range players {
      if p.ID == player_id {
        player = p
        goto Found
      }
    }
  }
  fmt.Println("Couldn't find player")
  return  // Didn't find the player
  Found:

  // Add to the team
  team.AddPlayer(player, true)

  // Remove from any position list
  for _, pos := range player.Positions {
    delete(sim.Avail_players[pos], player.ID)
  }

  sim.ResetStats()
}

func (sim *Simulation) RemoveKeeper(keeper_id string, team_id int) {
  // TODO
  sim.ResetStats()
}

// ResetStats resets every player's stats to zero
func (sim *Simulation) ResetStats() {
  for _, p := range sim.All_players {
    p.ResetWins()
  }

  // TODO: Also reset some stats for teams?
}
