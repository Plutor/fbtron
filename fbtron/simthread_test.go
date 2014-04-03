package fbtron

import (
  "testing"
)

func FakeSimulation() Simulation {
  var sim Simulation

  // Fake init players
  sim.Avail_players = BuildPlayersFromCsv("testdata/players_csv_ok.csv", "")

  // Init teams
  *num_teams = 2
  sim.InitTeams(map[string]int {
    "1B": 1,
    "SP": 1,
  })

  return sim
}

func TestRunSimulation(t *testing.T) {
  // TODO: Test starting a simulation thread, ask it for data, make sure data is
  // good.
}

func TestRunSeason(t *testing.T) {
  sim := FakeSimulation()

  sim.RunSeason()
  if v := sim.Num_seasons; v != 1 {
    t.Errorf("Error running season: expected 1, got %d", v)
  }
}

func TestInitPlayers(t *testing.T) {
  sim := FakeSimulation()

  // TODO: Fix no such file or directory error
  // sim.InitPlayers()
  if v := len(sim.Avail_players); v <= 0 {
    t.Errorf("InitPlayers: expected to load >0 players, got %d", v)
  }
}

func TestInitTeams(t *testing.T) {
  sim := FakeSimulation()
  // FakeSimulation() calls InitTeams() for us.

  if v := len(sim.Teams); v != 2 {
    t.Errorf("InitTeams: expected 2 teams created, got %d", v)
  }
  for n := range sim.Teams {
    if v := sim.Teams[n].GetOpenPosition(); v == "" {
      t.Errorf("InitTeams: expected non-blank empty position, got '%s'", v)
    }
  }
}

func TestDoDraft(t *testing.T) {
  sim := FakeSimulation()

  sim.DoDraft()
  for n := range sim.Teams {
    if v := sim.Teams[n].GetOpenPosition(); v != "" {
      t.Errorf("DoDraft: team %d expected blank empty position, got '%s'", n, v)
    }
    if len(sim.Avail_players) != 1 {
      t.Errorf("DoDraft: team %d expected 1 available player, got\n%s",
               n, sim.Avail_players)
    }
  }
}

func BenchmarkDoDraft(b *testing.B) {
  for i := 0; i < b.N; i++ {
    sim := FakeSimulation()
    sim.DoDraft()
  }
}

func TestRandomAvailablePlayerIndex(t *testing.T) {
  sim := FakeSimulation()

  // Remove players one by one, requesting a random player index, and making
  // sure it always falls in the range of 0<n<len(players)
  for ; len(sim.Avail_players) > 0;
      sim.Avail_players = sim.Avail_players[:len(sim.Avail_players)-1] {
    valid_pos := sim.Avail_players[0].positions[0]
    if v := sim.RandomAvailablePlayerIndex(valid_pos);
        v < 0 || v >= len(sim.Avail_players) {
      t.Errorf("RandomAvailablePlayerIndex: expected 0<index<%d, got %d",
               len(sim.Avail_players), v)
    }
  }
}

func TestAllAvailablePlayersIndexes(t *testing.T) {
  sim := FakeSimulation()

  if v := sim.AllAvailablePlayersIndexes("1B"); len(v) != 2 {
    t.Errorf("AllAvailablePlayersIndexes: expected len() == 2, got %d", v)
  } else {
    for n := range v {
      if v[n] < 0 || v[n] >= len(sim.Avail_players) {
        t.Errorf("AllAvailablePlayersIndexes: expected 0<index<%d, got %d",
                 len(sim.Avail_players), v[n])
      }
    }
  }

  if v := sim.AllAvailablePlayersIndexes("SP"); len(v) != 3 {
    t.Errorf("AllAvailablePlayersIndexes: expected len() == 3, got %d", v)
  } else {
    for n := range v {
      if v[n] < 0 || v[n] >= len(sim.Avail_players) {
        t.Errorf("AllAvailablePlayersIndexes: expected 0<index<%d, got %d",
                 len(sim.Avail_players), v[n])
      }
    }
  }
}

func TestScoreSeason(t *testing.T) {
  sim := FakeSimulation()

  var p *Player
  p = new(Player)
  p.positions = []string {"SP"}
  p.SetStat("R", 1.0)
  sim.Teams[0].AddPlayer(p, false)
  p = new(Player)
  p.positions = []string {"SP"}
  p.SetStat("R", 2.0)
  sim.Teams[1].AddPlayer(p, false)

  sim.ScoreSeason()

  if v:= sim.Teams[0].wins; v != 0 {
    t.Errorf("ScoreSeason: Expected team 0 wins == 0, got %d", v)
  }
  if v:= sim.Teams[1].wins; v != 1 {
    t.Errorf("ScoreSeason: Expected team 1 wins == 1, got %d", v)
  }
}

func TestEndSeason(t *testing.T) {
  sim := FakeSimulation()

  num_players := len(sim.Avail_players)
  sim.DoDraft()

  if v := len(sim.Avail_players); v >= num_players {
    t.Errorf("EndSeason: expected numplayers decreased, got %d >= %d",
             v, num_players)
  }

  sim.EndSeason()
  if v := len(sim.Avail_players); v != num_players {
    t.Errorf("EndSeason: expected numplayers unchanged, got %d != %d",
             v, num_players)
  }

  // TODO: Also add wins and make sure the players all got the wins applied.
}
