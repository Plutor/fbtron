package fbtron

import (
  "sort"
  "strconv"
  "testing"
)

// FakeTeam generates a fake team for testing. It has 6 players on its roster,
// each named after their index. Every third one (starting with the first) is a
// keeper (so 2 total). The team has 10 wins.
func FakeTeam() *Team {
  team := new(Team)
  team.wins = 10
  team.SetPositions(map[string]int {
    "1B": 1,
    "SP": 1,
    "Fake": 6,
  })

  for n := 0; n < 6; n++ {
    player := new(Player)
    player.Firstname = strconv.Itoa(n)
    player.Positions = []string { "Fake" }
    team.AddPlayer(player, n % 3 == 0, "Fake")
  }

  return team
}

func TestGetAllOpenPositions(t *testing.T) {
  team := FakeTeam()

  v1 := team.GetAllOpenPositions()
  sort.Strings(v1)
  if len(v1) != 2 || v1[0] != "1B" || v1[1] != "SP" {
    t.Errorf("GetOpenPosition: expected [1B, SP], got '%s'", v1)
  }

  player := Player {
      Firstname: "Openposition",
      Lastname: "Filler",
      Positions: []string { v1[0] },
  }
  team.AddPlayer(&player, false, v1[0])

  v2 := team.GetAllOpenPositions()
  if len(v2) != 1 && v2[0] != "SP" {
    t.Errorf("GetOpenPosition: expected [SP], got '%s'", v1)
  }
}


func TestTeamAddPlayer(t *testing.T) {
  team := FakeTeam()
  if v := len(team.Roster["Fake"]); v != 6 {
    t.Errorf("Error adding players to a team: expected 6 Fakes, got %d", v)
    t.FailNow()
  }

  for n := range team.Roster["Fake"] {
    if v := team.Roster["Fake"][n].Player.GetName(); v != strconv.Itoa(n) {
      t.Errorf("Error adding player to a team: player %d " +
               "expected name '%s', got '%s'",
               n, strconv.Itoa(n), v )
      t.FailNow()
    }
  }

  player := Player {
      Firstname: "Openposition",
      Lastname: "Filler",
      Positions: []string { "1B", "SP" },
  }
  team.AddPlayer(&player, false, "1B")
  if len(team.Roster["1B"]) != 1 {
    t.Errorf("Error adding 1B to a team", team.Roster)
  }

  player = Player {
      Firstname: "Openposition",
      Lastname: "Filler, Jr.",
      Positions: []string { "1B", "SP" },
  }
  team.AddPlayer(&player, false, "SP")
  if len(team.Roster["SP"]) != 1 {
    t.Errorf("Error adding a SP to a team")
  }
}

func TestRelease(t *testing.T) {
  team := FakeTeam()

  released := team.Release()
  if len(released) != 4 {
    t.Errorf("Error releasing non-keeper players: " +
             "expected 4 released, got %d", len(released))
  }
  if v := len(team.Roster["Fake"]); v != 2 {
    t.Errorf("Error releasing non-keeper players: " +
             "expected 2 remaining, got %d", v)
  }
}

func TestCreditRosterWithWins(t *testing.T) {
  team := FakeTeam()
  team.CreditRosterWithWins()

  for n := range team.Roster["Fake"] {
    if v := team.Roster["Fake"][n].Player.Total_wins; v != team.wins {
      t.Errorf("Error crediting roster: team has %d wins, player '%s' has %d",
               team.wins, team.Roster["Fake"][n].Player.GetName(), v)
    }
    if v := team.Roster["Fake"][n].Player.Num_seasons; v != 1 {
      t.Errorf("Error crediting roster: expected Num_seasons=1, " +
               "player '%s' has %d",
               team.Roster["Fake"][n].Player.GetName(), v)
    }
  }
}

func TestGetTeamStat(t *testing.T) {
  team := FakeTeam()

  // Test a summed stat
  team.Roster["Fake"][0].Player.SetStat("B_R", 1)
  team.Roster["Fake"][1].Player.SetStat("B_R", 1)
  team.Roster["Fake"][2].Player.SetStat("B_R", 40)
  if v := team.GetStat("B_R"); v != 42 {
    t.Errorf("Error with summed stat, expected 42, got %f", v)
  }

  // Test an ab-weighted stat
  team.Roster["Fake"][0].Player.SetStat("B_AVG", 0.200)
  team.Roster["Fake"][0].Player.SetStat("B_AB", 10)
  team.Roster["Fake"][1].Player.SetStat("B_AVG", 0.200)
  team.Roster["Fake"][1].Player.SetStat("B_AB", 10)
  team.Roster["Fake"][2].Player.SetStat("B_AVG", 0.500)
  team.Roster["Fake"][2].Player.SetStat("B_AB", 20)
  if v := team.GetStat("B_AVG"); v != 0.350 {
    t.Errorf("Error with ab-weighted stat, expected 0.350, got %f", v)
  }

  // Test an ip-weighted stat
  team.Roster["Fake"][0].Player.SetStat("P_ERA", 2.00)
  team.Roster["Fake"][0].Player.SetStat("P_IP", 10)
  team.Roster["Fake"][1].Player.SetStat("P_ERA", 2.00)
  team.Roster["Fake"][1].Player.SetStat("P_IP", 10)
  team.Roster["Fake"][2].Player.SetStat("P_ERA", 5.00)
  team.Roster["Fake"][2].Player.SetStat("P_IP", 20)
  if v := team.GetStat("P_ERA"); v != -3.50 {
    t.Errorf("Error with descending ip-weighted stat, expected -3.50, got %f",
             v)
  }

  // Test an unknown stat
  if v := team.GetStat("ZOMGBBQ"); v != 0.0 {
    t.Errorf("Error with unknown stat, expected 0, got %f", v)
  }
}

//
// Benchmarks
//

func BenchmarkGetAllOpenPositions(b *testing.B) {
  team := FakeTeam()
  team.SetPositions(POSITIONS)

  b.ResetTimer()
  for n := 0; n < b.N; n++ {
    team.GetAllOpenPositions()
  }
}

func BenchmarkAddPlayer(b *testing.B) {
  teams := make([]Team, b.N)
  players := []Player {
      Player {
        Firstname: "Jimmy",
        Lastname: "Firstbaseman",
        Positions: []string { "1B" },
      },
      Player {
        Firstname: "Joey",
        Lastname: "Secondbaseman",
        Positions: []string { "2B" },
      },
  }

  for n := range teams {
    teams[n] = Team {}
    teams[n].SetPositions(map[string]int {
      "1B": 1,
      "2B": 1,
    })
  }

  b.ResetTimer()
  for n := 0; n < b.N; n++ {
    teams[n].AddPlayer(&players[0], false, "1B")
    teams[n].AddPlayer(&players[1], false, "2B")
  }
}

