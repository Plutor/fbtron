package fbtron

import (
  "testing"
)

func FakeSimulation() Simulation {
  var sim Simulation

  // Fake init players
  sim.Avail_players = BuildPlayersFromCsv("testdata/players_csv_ok.csv", "X")

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

  // TODO
}

func TestDoDraft(t *testing.T) {
  // sim := FakeSimulation()

  // TODO
}

func TestScoreSeason(t *testing.T) {
  // sim := FakeSimulation()

  // TODO
}

func TestEndSeason(t *testing.T) {
  // sim := FakeSimulation()

  // TODO
}
