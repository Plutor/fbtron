package fbtron

type TeamMember struct {
  player    *Player
  keeper    bool
}

type Team struct {
  name        string
  roster      map[string][]TeamMember
  wins        int
}

func (team *Team) SetPositions(positions map[string]int) {
  team.roster = make(map[string][]TeamMember)
  for position, num := range positions {
    team.roster[position] = make([]TeamMember, 0, num)
  }
}

// Finds a position with an open spot on the roster and returns it. Returns
// empty string if there are no open positions.
func (team *Team) GetOpenPosition() string {
  for position, members := range team.roster {
    if cap(members) > len(members) {
      return position
    }
  }
  return ""
}

// AddPlayer adds the passed player to the team roster. If keeper is true, the
// player will not be released by the Release() call.
func (team *Team) AddPlayer(p *Player, keeper bool) {
  // TODO: Select the position properly for players with multiple positions.
  pos := p.positions[0]
  if len(team.roster[pos]) == cap(team.roster[pos]) {
    // TODO: This is a bad problem -- prevent it!
    // No open spots! ONO
    return
  }

  // Grow the array by one and add this player
  team.roster[pos] = team.roster[pos][:len(team.roster[pos])+1]
  team.roster[pos][len(team.roster[pos])-1] = TeamMember {p, keeper}
}

// Release releases all of the players in the roster that are not marked as
// keepers. Returns an array of the released players.
func (team *Team) Release() []*Player {
  team.CreditRosterWithWins()

  released := make([]*Player, 0)
  for position, members := range team.roster {
    newmembers := make([]TeamMember, 0, len(members))
    for n := range members {
      if members[n].keeper {
        newmembers = newmembers[:len(newmembers)+1]
        newmembers[len(newmembers)-1] = members[n]
      } else {
        released = append(released, members[n].player)
      }
    }
    team.roster[position] = newmembers
  }

  return released
}

// CreditRosterWithWins adds the wins for this team to every player on the
// roster, and also increments the number of seasons. This should only be called
// once per season, ideally by Release().
func (team *Team) CreditRosterWithWins() {
  for _, members := range team.roster {
    for n := range members {
      members[n].player.total_wins += team.wins
      members[n].player.num_seasons++
    }
  }
}

// GetStat gets the team-wide value for a stat, either summed or ip/ab-weighted,
// depending on the stat.
func (team *Team) GetStat(statname string) float64 {
  switch GetStatType(statname) {
  case STAT_SUMMED:
    sum := 0.0
    for _, members := range team.roster {
      for n := range members {
        sum += members[n].player.GetStat(statname)
      }
    }
    return sum
  case STAT_AB_WEIGHTED_AVG:
    avg := 0.0
    ab := 0.0
    for _, members := range team.roster {
      for n := range members {
        p := members[n].player
        avg += p.GetStat(statname) * p.GetStat("AB")
        ab += p.GetStat("AB")
      }
    }
    return avg / ab
  case STAT_IP_WEIGHTED_AVG:
    avg := 0.0
    ip := 0.0
    for _, members := range team.roster {
      for n := range members {
        p := members[n].player
        avg += p.GetStat(statname) * p.GetStat("IP")
        ip += p.GetStat("IP")
      }
    }
    return avg / ip
  }

  return 0.0
}
