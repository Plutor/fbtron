package fbtron

type TeamMember struct {
  Player    *Player
  Keeper    bool
}

type Team struct {
  Name        string
  Roster      []TeamMember
  positions   []string
  wins        int
}

func (team *Team) SetPositions(positions map[string]int) {
  team.positions = make([]string, 0, len(positions))
  for pos, num := range positions {
    for i := 0; i < num; i++ {
      team.positions = append(team.positions, pos)
    }
  }
  team.Roster = make([]TeamMember, 0, len(team.positions))
}

// Finds a position with an open spot on the roster and returns it. Returns
// empty string if there are no open positions.
func (team *Team) GetAllOpenPositions() []string {
  positions := team.positions

  for n := range team.Roster {
    Nextplayer:
    for _, pos := range team.Roster[n].Player.Positions {
      for n, avail_pos := range positions {
        if pos == avail_pos {
          positions = append(positions[:n], positions[n+1:]...)
          break Nextplayer
        }
      }
    }
  }
  return positions
}

// AddPlayer adds the passed player to the team roster. If keeper is true, the
// player will not be released by the Release() call.
func (team *Team) AddPlayer(p *Player, keeper bool) {
  team.Roster = append(team.Roster, TeamMember {p, keeper})
}

// Release releases all of the players in the roster that are not marked as
// keepers. Returns an array of the released players.
func (team *Team) Release() []*Player {
  team.CreditRosterWithWins()

  released := make([]*Player, 0, len(team.Roster))
  newmembers := make([]TeamMember, 0, len(team.Roster))
  for n := range team.Roster {
    if team.Roster[n].Keeper {
      newmembers = append(newmembers, team.Roster[n])
    } else {
      released = append(released, team.Roster[n].Player)
    }
  }
  team.Roster = newmembers

  return released
}

// HasPlayer returns true if this player is on this team's roster. If
// keeper_only is true, returns true only if the player is also a keeper.
func (team *Team) HasPlayer(player_id string, keeper_only bool) bool {
  for n := range team.Roster {
    if team.Roster[n].Player != nil && team.Roster[n].Player.ID == player_id {
      return !keeper_only || team.Roster[n].Keeper
    }
  }

  return false
}

// CreditRosterWithWins adds the wins for this team to every player on the
// roster, and also increments the number of seasons. This should only be called
// once per season, ideally by Release().
func (team *Team) CreditRosterWithWins() {
  for n := range team.Roster {
    team.Roster[n].Player.Total_wins += team.wins
    team.Roster[n].Player.Num_seasons++
  }
}

// GetStat gets the team-wide value for a stat, either summed or ip/ab-weighted,
// depending on the stat. Returns a negative value for stats that are sorted
// descending.
func (team *Team) GetStat(statname string) float64 {
  rv := 0.0
  st := GetStatType(statname)
  switch {
  case st & STAT_SUMMED != 0:
    for n := range team.Roster {
      rv += team.Roster[n].Player.GetStat(statname)
    }
  case st & STAT_AB_WEIGHTED_AVG != 0:
    avg := 0.0
    ab := 0.0
    for n := range team.Roster {
      p := team.Roster[n].Player
      avg += p.GetStat(statname) * p.GetStat("B_AB")
      ab += p.GetStat("B_AB")
    }
    rv = avg / ab
  case st & STAT_IP_WEIGHTED_AVG != 0:
    avg := 0.0
    ip := 0.0
    for n := range team.Roster {
      p := team.Roster[n].Player
      avg += p.GetStat(statname) * p.GetStat("P_IP")
      ip += p.GetStat("P_IP")
    }
    rv = avg / ip
  }

  if st & STAT_DESC != 0 {
    return -rv
  }
  return rv
}
