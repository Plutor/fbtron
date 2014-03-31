package fbtron

import (
  "strconv"
  "testing"
)

func TestTeamAddPlayer(t *testing.T) {
  team := new(Team)
  for n := 0; n < 99; n++ {
    player := new(Player)
    player.name = strconv.Itoa(n)
    team.AddPlayer(player, true)
    if v := len(team.roster); v != n + 1 {
      t.Errorf("Error adding player to a team: expected %d, got %d", n + 1, v)
      t.FailNow()
    }
  }

  for n := 0; n < 99; n++ {
    if team.roster[n].name != strconv.Itoa(n) {
      t.Errorf("Error adding player to a team: player %d " +
               "expected name '%s', got '%s'",
               n, strconv.Itoa(n), team.roster[n].name )
      t.FailNow()
    }
  }
}

func TestRelease(t *testing.T) {
  team := new(Team)
  for n := 0; n < 99; n++ {
    player := new(Player)
    team.AddPlayer(player, n < 50)
    if v := len(team.roster); v != n + 1 {
      t.Errorf("Error adding player to a team: expected %d, got %d", n + 1, v)
    }
  }

  released := team.Release()
  if len(released) != 49 {
    t.Errorf("Error releasing non-keeper players: " +
             "expected 49 released, got %d", len(released))
  }
  if v := len(team.roster); v != 50 {
    t.Errorf("Error releasing non-keeper players: " +
             "expected 50 remaining, got %d", v)
  }
}

// TODO: Get team-wide stats
