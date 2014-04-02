package fbtron

import (
  "strconv"
  "testing"
)

// FakeTeam generates a fake team for testing. It has 99 players on its roster,
// each named after their index. Every third one (starting with the first) is a
// keeper (so 33 total). The team has 10 wins.
func FakeTeam() *Team {
  team := new(Team)
  team.wins = 10
  team.SetPositions(map[string]int {
    "1B": 1,
    "SP": 1,
    "Fake": 99,
  })

  for n := 0; n < 99; n++ {
    player := new(Player)
    player.firstname = strconv.Itoa(n)
    player.positions = []string { "Fake" }
    team.AddPlayer(player, n % 3 == 0)
  }

  return team
}

func TestGetOpenPosition(t *testing.T) {
  team := FakeTeam()

  v1 := team.GetOpenPosition()
  if v1 != "1B" && v1 != "SP" {
    t.Errorf("GetOpenPosition: expected 1B or SP, got '%s'", v1)
  }

  player := Player {
      firstname: "Openposition",
      lastname: "Filler",
      positions: []string { v1 },
  }
  team.AddPlayer(&player, false)

  v2 := team.GetOpenPosition()
  if v2 != "1B" && v2 != "SP" {
    t.Errorf("GetOpenPosition: expected 1B or SP, got '%s'", v1)
  } else if v1 == v2 {
    t.Errorf("GetOpenPosition: expected '%s' != '%s'", v1, v2)
  }
}


func TestTeamAddPlayer(t *testing.T) {
  team := FakeTeam()
  if v := len(team.roster["Fake"]); v != 99 {
    t.Errorf("Error adding players to a team: expected 99 Fakes, got %d", v)
    t.FailNow()
  }

  for n := range team.roster["Fake"] {
    if v := team.roster["Fake"][n].player.GetName(); v != strconv.Itoa(n) {
      t.Errorf("Error adding player to a team: player %d " +
               "expected name '%s', got '%s'",
               n, strconv.Itoa(n), v )
      t.FailNow()
    }
  }
}

func TestRelease(t *testing.T) {
  team := FakeTeam()

  released := team.Release()
  if len(released) != 66 {
    t.Errorf("Error releasing non-keeper players: " +
             "expected 66 released, got %d", len(released))
  }
  if v := len(team.roster["Fake"]); v != 33 {
    t.Errorf("Error releasing non-keeper players: " +
             "expected 33 remaining, got %d", v)
  }
}

func TestCreditRosterWithWins(t *testing.T) {
  team := FakeTeam()
  team.CreditRosterWithWins()

  for n := range team.roster["Fake"] {
    if v := team.roster["Fake"][n].player.total_wins; v != team.wins {
      t.Errorf("Error crediting roster: team has %d wins, player '%s' has %d",
               team.wins, team.roster["Fake"][n].player.GetName(), v)
    }
    if v := team.roster["Fake"][n].player.num_seasons; v != 1 {
      t.Errorf("Error crediting roster: expected num_seasons=1, " +
               "player '%s' has %d",
               team.roster["Fake"][n].player.GetName(), v)
    }
  }
}

func TestGetTeamStat(t *testing.T) {
  team := FakeTeam()

  // Test a summed stat
  team.roster["Fake"][0].player.SetStat("R", 1)
  team.roster["Fake"][1].player.SetStat("R", 2)
  team.roster["Fake"][2].player.SetStat("R", 3)
  team.roster["Fake"][3].player.SetStat("R", 10)
  team.roster["Fake"][4].player.SetStat("R", 12)
  team.roster["Fake"][5].player.SetStat("R", 14)
  if v := team.GetStat("R"); v != 42 {
    t.Errorf("Error with summed stat, expected 42, got %f", v)
  }

  // Test an ab-weighted stat
  team.roster["Fake"][6].player.SetStat("BA", 0.200)
  team.roster["Fake"][6].player.SetStat("AB", 10)
  team.roster["Fake"][7].player.SetStat("BA", 0.200)
  team.roster["Fake"][7].player.SetStat("AB", 10)
  team.roster["Fake"][8].player.SetStat("BA", 0.500)
  team.roster["Fake"][8].player.SetStat("AB", 20)
  if v := team.GetStat("BA"); v != 0.350 {
    t.Errorf("Error with ab-weighted stat, expected 0.350, got %f", v)
  }

  // Test an ip-weighted stat
  team.roster["Fake"][6].player.SetStat("ERA", 2.00)
  team.roster["Fake"][6].player.SetStat("IP", 10)
  team.roster["Fake"][7].player.SetStat("ERA", 2.00)
  team.roster["Fake"][7].player.SetStat("IP", 10)
  team.roster["Fake"][8].player.SetStat("ERA", 5.00)
  team.roster["Fake"][8].player.SetStat("IP", 20)
  if v := team.GetStat("ERA"); v != 3.50 {
    t.Errorf("Error with ip-weighted stat, expected 3.50, got %f", v)
  }

  // Test an unknown stat
  if v := team.GetStat("ZOMGBBQ"); v != 0.0 {
    t.Errorf("Error with unknown stat, expected 0, got %f", v)
  }

}


// TODO: Get team-wide stats
