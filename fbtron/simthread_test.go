package fbtron

import (
  "fmt"
  "testing"
)

func FakeSimulation(csv_suffix string) Simulation {
  var sim Simulation

  // Fake init players
  sim.AddPlayersToPositionLists(BuildPlayersFromCsv(
      fmt.Sprintf("testdata/players_csv_ok%s.csv", csv_suffix), ""))

  // Init teams
  *num_teams = 10
  sim.InitTeams(map[string]int {
    "1B": 1,
    "2B": 1,
  })

  return sim
}

func TestRunSimulation(t *testing.T) {
  // TODO: Test starting a simulation thread, ask it for data, make sure data is
  // good.
}

func TestRunSeason(t *testing.T) {
  sim := FakeSimulation("")

  sim.RunSeason()
  if v := sim.Num_seasons; v != 1 {
    t.Errorf("Error running season: expected 1, got %d", v)
  }
}

func TestInitPlayers(t *testing.T) {
  sim := FakeSimulation("")

  // TODO: Fix no such file or directory error
  // sim.InitPlayers()
  if v := len(sim.Avail_players); v <= 0 {
    t.Errorf("InitPlayers: expected to load >0 positions, got %d", v)
  }
  if v := len(sim.Avail_players["2B"]); v != 3 {
    t.Errorf("InitPlayers: expected to load 3 2B, got %d", v)
  }
  if v := len(sim.Avail_players["1B"]); v != 2 {
    t.Errorf("InitPlayers: expected to load 2 1B, got %d", v)
  }
}

func TestInitTeams(t *testing.T) {
  sim := FakeSimulation("")
  // FakeSimulation() calls InitTeams() for us.

  if v := len(sim.Teams); v != 10 {
    t.Errorf("InitTeams: expected 10 teams created, got %d", v)
  }
  for n := range sim.Teams {
    if v := sim.Teams[n].GetOpenPosition(); v == "" {
      t.Errorf("InitTeams: expected non-blank empty position, got '%s'", v)
    }
  }
}

func TestDoDraft(t *testing.T) {
  sim := FakeSimulation("_big")
  sim.DoDraft()
  for n := range sim.Teams {
    if v := sim.Teams[n].GetOpenPosition(); v != "" {
      t.Errorf("DoDraft: team %d expected blank empty position, got '%s'", n, v)
    }
  }

  num_players := 0
  for _, players := range sim.Avail_players {
    num_players += len(players)
  }
  if num_players == 0 {
    t.Errorf("DoDraft: expected >0 available player after draft, got %d",
             num_players)
  }
}

func BenchmarkDoDraft(b *testing.B) {
  orig_sim := FakeSimulation("_big")

  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    sim := orig_sim
    sim.DoDraft()
  }
}

func TestAddPlayersToPositionLists(t *testing.T) {
  sim := FakeSimulation("")

  if v := len(sim.Avail_players); v != 2 {
    t.Errorf("AddPlayersToPositionLists: expected 2 position indexes, got %d",
             v)
  }
  if v := len(sim.Avail_players["2B"]); v != 3 {
    t.Errorf("AddPlayersToPositionLists: expected 3 2B, got %d", v)
  }
  if v := len(sim.Avail_players["1B"]); v != 2 {
    t.Errorf("AddPlayersToPositionLists: expected 2 1B, got %d", v)
  }
}

func TestRandomAvailablePlayer(t *testing.T) {
  sim := FakeSimulation("")

  // Remove players one by one, requesting a random player index, and making
  // sure it always falls in the range of 0<n<len(players)
  for _, pos := range []string{"2B", "1B"} {
    for len(sim.Avail_players[pos]) > 0 {
      if v := sim.RandomAvailablePlayer(pos); v == nil {
        t.Errorf("RandomAvailablePlayer: expected a Player for %s, got %s",
                 pos, v)
      }
    }
  }
}

func TestScoreSeason(t *testing.T) {
  sim := FakeSimulation("")

  var p *Player
  p = new(Player)
  p.Positions = []string {"2B"}
  p.SetStat("B_R", 1.0)
  sim.Teams[0].AddPlayer(p, false)
  p = new(Player)
  p.Positions = []string {"2B"}
  p.SetStat("B_R", 2.0)
  sim.Teams[1].AddPlayer(p, false)

  sim.ScoreSeason()

  if v:= sim.Teams[0].wins; v != 8 {
    t.Errorf("ScoreSeason: Expected team 0 wins == 8, got %d", v)
  }
  if v:= sim.Teams[1].wins; v != 9 {
    t.Errorf("ScoreSeason: Expected team 1 wins == 9, got %d", v)
  }
}

func BenchmarkScoreSeason(b *testing.B) {
  orig_sim := FakeSimulation("_big")
  orig_sim.DoDraft()
  b.ResetTimer()

  for i := 0; i < b.N; i++ {
    sim := orig_sim
    sim.ScoreSeason()
  }
}

func TestEndSeason(t *testing.T) {
  sim := FakeSimulation("")

  num_players_start := 0
  for _, players := range sim.Avail_players {
    num_players_start += len(players)
  }

  sim.DoDraft()

  // End the season and expect the number of players to be back where we
  // started.
  sim.EndSeason()

  num_players := 0
  for _, players := range sim.Avail_players {
    num_players += len(players)
  }
  if num_players != num_players_start {
    t.Errorf("EndSeason: expected num_players unchanged, got %d != %d",
             num_players, num_players_start)
  }

  // TODO: Also add wins and make sure the players all got the wins applied.
}

func BenchmarkEndSeason(b *testing.B) {
  orig_sim := FakeSimulation("_big")
  orig_sim.DoDraft()
  b.ResetTimer()

  for i := 0; i < b.N; i++ {
    sim := orig_sim
    sim.EndSeason()
  }
}

func TestMerge(t *testing.T) {
  // TODO
}

func BenchmarkRandomAvailablePlayer(b *testing.B) {
  sim := FakeSimulation("_big")
  for i := 0; i < b.N; i++ {
    sim.RandomAvailablePlayer("1B")
  }
}
