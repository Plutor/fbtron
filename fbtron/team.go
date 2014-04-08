package fbtron

type TeamMember struct {
  Player    *Player
  Keeper    bool
}

type Team struct {
  Name        string
  Roster      map[string][]TeamMember
  wins        int
}

func (team *Team) SetPositions(positions map[string]int) {
  team.Roster = make(map[string][]TeamMember)
  for position, num := range positions {
    team.Roster[position] = make([]TeamMember, 0, num)
  }
}

// Finds a position with an open spot on the roster and returns it. Returns
// empty string if there are no open positions.
func (team *Team) GetOpenPosition() string {
  for position, members := range team.Roster {
    if cap(members) > len(members) {
      return position
    }
  }
  return ""
}

// AddPlayer adds the passed player to the team roster. If keeper is true, the
// player will not be released by the Release() call.
func (team *Team) AddPlayer(p *Player, keeper bool) {
  // Select an open position for players with multiple positions.
  var pos string
  for _, playerpos := range p.Positions {
    if len(team.Roster[playerpos]) < cap(team.Roster[playerpos]) {
      pos = playerpos
    }
  }
  if pos == "" {
    // TODO: No open positions is potentially a bad problem. Instead, create a
    // special "overflow" player type that is used in this case.
    return
  }

  // Grow the array by one and add this player
  team.Roster[pos] = team.Roster[pos][:len(team.Roster[pos])+1]
  team.Roster[pos][len(team.Roster[pos])-1] = TeamMember {p, keeper}
}

// Release releases all of the players in the roster that are not marked as
// keepers. Returns an array of the released players.
func (team *Team) Release() []*Player {
  team.CreditRosterWithWins()

  released := make([]*Player, 0)
  for position, members := range team.Roster {
    newmembers := make([]TeamMember, 0, len(members))
    for n := range members {
      if members[n].Keeper {
        newmembers = newmembers[:len(newmembers)+1]
        newmembers[len(newmembers)-1] = members[n]
      } else {
        released = append(released, members[n].Player)
      }
    }
    team.Roster[position] = newmembers
  }

  return released
}

// CreditRosterWithWins adds the wins for this team to every player on the
// roster, and also increments the number of seasons. This should only be called
// once per season, ideally by Release().
func (team *Team) CreditRosterWithWins() {
  for _, members := range team.Roster {
    for n := range members {
      members[n].Player.Total_wins += team.wins
      members[n].Player.Num_seasons++
    }
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
    for _, members := range team.Roster {
      for n := range members {
        rv += members[n].Player.GetStat(statname)
      }
    }
  case st & STAT_AB_WEIGHTED_AVG != 0:
    avg := 0.0
    ab := 0.0
    for _, members := range team.Roster {
      for n := range members {
        p := members[n].Player
        avg += p.GetStat(statname) * p.GetStat("B_AB")
        ab += p.GetStat("B_AB")
      }
    }
    rv = avg / ab
  case st & STAT_IP_WEIGHTED_AVG != 0:
    avg := 0.0
    ip := 0.0
    for _, members := range team.Roster {
      for n := range members {
        p := members[n].Player
        avg += p.GetStat(statname) * p.GetStat("P_IP")
        ip += p.GetStat("P_IP")
      }
    }
    rv = avg / ip
  }

  if st & STAT_DESC != 0 {
    return -rv
  }
  return rv
}
