package fbtron

import (
  "testing"
)

// TODO: Separate roster into actually drafted and simulated draft

func TestTeamAddPlayer(t *testing.T) {
  team := new(Team)
  for n := 1; n < 99; n++ {
    player := new(Player)
    team.AddPlayer(player)
    if v := len(team.roster); v != n {
      t.Errorf("Error adding player to a team: expected %d, got %d", n, v)
    }
  }
}

// TODO: Release (releases simulated only)

// TODO: Get team-wide stats
