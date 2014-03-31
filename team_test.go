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

  for n := 0; n < 99; n++ {
    player := new(Player)
    player.name = strconv.Itoa(n)
    team.AddPlayer(player, n % 3 == 0)
  }

  return team
}

func TestTeamAddPlayer(t *testing.T) {
  team := FakeTeam()
  if v := len(team.roster); v != 99 {
    t.Errorf("Error adding player to a team: expected 99, got %d", v)
    t.FailNow()
  }

  for n := range team.roster {
    if team.roster[n].name != strconv.Itoa(n) {
      t.Errorf("Error adding player to a team: player %d " +
               "expected name '%s', got '%s'",
               n, strconv.Itoa(n), team.roster[n].name )
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
  if v := len(team.roster); v != 33 {
    t.Errorf("Error releasing non-keeper players: " +
             "expected 33 remaining, got %d", v)
  }
}

func TestCreditRosterWithWins(t *testing.T) {
  team := FakeTeam()
  team.CreditRosterWithWins()

  for n := range team.roster {
    if v := team.roster[n].total_wins; v != team.wins {
      t.Errorf("Error crediting roster: team has %d wins, player '%s' has %d",
               team.wins, team.roster[n].name, v)
    }
    if v := team.roster[n].num_seasons; v != 1 {
      t.Errorf("Error crediting roster: expected num_seasons=1, " +
               "player '%s' has %d",
               team.roster[n].name, v)
    }
  }
}


// TODO: Get team-wide stats
