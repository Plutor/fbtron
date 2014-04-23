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

// Tests and benchmarks for RunSimulation are in main_test.go

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

  num_after_season := 0
  for _, players := range sim.Avail_players {
    num_after_season += len(players)
  }
  if num_after_season == 0 {
    t.Errorf("DoDraft: expected >0 available player after draft, got %d",
             num_after_season)
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

func TestEndSeason(t *testing.T) {
  sim := FakeSimulation("")

  num_before_draft := 0
  for _, players := range sim.Avail_players {
    num_before_draft += len(players)
  }

  sim.DoDraft()

  num_after_draft := 0
  for _, players := range sim.Avail_players {
    num_after_draft += len(players)
  }

  for n := range sim.Teams {
    sim.Teams[n].wins = n+1
  }

  // End the season and expect the number of players to be back where we
  // started.
  sim.EndSeason()

  // Make sure the players all got released and got their wins applied.
  num_after_season := 0
  num_with_wins := 0
  num_with_seasons := 0
  for _, players := range sim.Avail_players {
    num_after_season += len(players)
    for _, p := range players {
      if p.Total_wins > 0 {
        num_with_wins++
      }
      if p.Num_seasons > 0 {
        num_with_seasons++
      }
    }
  }

  if num_after_draft == num_before_draft {
    t.Errorf("EndSeason: expected before and after draft changed, got %d = %d",
             num_before_draft, num_after_draft)
  }
  if num_after_season != num_before_draft {
    t.Errorf("EndSeason: expected before and after season unchanged, " +
             "got %d != %d", num_before_draft, num_after_season)
  }
  if num_with_wins != num_before_draft - num_after_draft {
    t.Errorf("EndSeason: Not enough players got wins: expected %d, got %d",
             num_before_draft - num_after_draft, num_with_wins)
  }
  if num_with_seasons != num_before_draft - num_after_draft {
    t.Errorf("EndSeason: Not enough players got seasons: expected %d, got %d",
             num_before_draft - num_after_draft, num_with_seasons)
  }
}

func TestMerge(t *testing.T) {
  // TODO
}

func TestTopPlayers(t *testing.T) {
  // TODO
}

func TestAddPlayer(t *testing.T) {
  // TODO
}

func TestDeletePlayer(t *testing.T) {
  // TODO
}

func TestResetStats(t *testing.T) {
  // TODO
}

//
// Benchmarks
//

func BenchmarkRunSeason(b *testing.B) {
  sim := FakeSimulation("_big")

  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    sim.RunSeason()
  }
}

func BenchmarkDoDraft(b *testing.B) {
  sim := FakeSimulation("_big")

  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    sim.DoDraft()
    b.StopTimer()
    sim.ScoreSeason()   // untimed
    sim.EndSeason()     // untimed
    b.StartTimer()
  }
}

func BenchmarkScoreSeason(b *testing.B) {
  sim := FakeSimulation("_big")

  b.ResetTimer()
  b.StopTimer()
  for i := 0; i < b.N; i++ {
    sim.DoDraft()     // untimed
    b.StartTimer()
    sim.ScoreSeason()
    b.StopTimer()
    sim.EndSeason()   // untimed
  }
}

func BenchmarkEndSeason(b *testing.B) {
  sim := FakeSimulation("_big")

  b.ResetTimer()
  b.StopTimer()
  for i := 0; i < b.N; i++ {
    sim.DoDraft()       // untimed
    sim.ScoreSeason()   // untimed
    b.StartTimer()
    sim.EndSeason()
    b.StopTimer()
  }
}

func BenchmarkRandomAvailablePlayer(b *testing.B) {
  sim := FakeSimulation("_big")
  for i := 0; i < b.N; i++ {
    sim.RandomAvailablePlayer("1B")
  }
}
