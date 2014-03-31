package fbtron

type Team struct {
  name        string
  roster      []*Player
  keeper      []bool
  wins        int
}

// AddPlayer adds the passed player to the team roster. If keeper is true, the
// player will not be released by the Release() call.
func (team *Team) AddPlayer(p *Player, keeper bool) {
  // Grow the arrays if necessary. Always assume that they're the same length.
  nthplayer := len(team.roster) + 1
  if nthplayer > cap(team.roster) {
    newroster := make([]*Player, len(team.roster), nthplayer * 2)
    copy(newroster, team.roster)
    team.roster = newroster

    newkeeper := make([]bool, len(team.roster), nthplayer * 2)
    copy(newkeeper, team.keeper)
    team.keeper = newkeeper
  }

  // Add to the end of the arrays
  team.roster = team.roster[0:nthplayer]
  team.roster[nthplayer-1] = p
  team.keeper = team.keeper[0:nthplayer]
  team.keeper[nthplayer-1] = keeper
}

// Release releases all of the players in the roster that are not marked as
// keepers. Returns an array of the released players.
func (team *Team) Release() []*Player {
  team.CreditRosterWithWins()

  rostersize := len(team.roster)
  newroster := make([]*Player, 0, rostersize)
  newkeeper := make([]bool, 0, rostersize)
  released := make([]*Player, 0, rostersize)

  for i := 0; i < rostersize; i++ {
    if team.keeper[i] {
      newroster = newroster[:len(newroster)+1]
      newroster[len(newroster)-1] = team.roster[i]
      newkeeper = newkeeper[:len(newkeeper)+1]
      newkeeper[len(newkeeper)-1] = true
    } else {
      released = released[:len(released)+1]
      released[len(released)-1] = team.roster[i]
    }
  }

  team.roster = newroster
  team.keeper = newkeeper

  return released
}

// CreditRosterWithWins adds the wins for this team to every player on the
// roster, and also increments the number of seasons. This should only be called
// once per season, ideally by Release().
func (team *Team) CreditRosterWithWins() {
  for i := 0; i < len(team.roster); i++ {
    team.roster[i].total_wins += team.wins
    team.roster[i].num_seasons++
  }
}

// GetStat gets the team-wide value for a stat, either summed or ip/ab-weighted,
// depending on the stat.
func (team *Team) GetStat(statname string) float64 {
  switch GetStatType(statname) {
  case STAT_SUMMED:
    sum := 0.0
    for n := range team.roster {
      sum += team.roster[n].GetStat(statname)
    }
    return sum
  case STAT_AB_WEIGHTED_AVG:
    avg := 0.0
    ab := 0.0
    for n := range team.roster {
      avg += team.roster[n].GetStat(statname) * team.roster[n].GetStat("AB")
      ab += team.roster[n].GetStat("AB")
    }
    return avg / ab
  case STAT_IP_WEIGHTED_AVG:
    avg := 0.0
    ip := 0.0
    for n := range team.roster {
      avg += team.roster[n].GetStat(statname) * team.roster[n].GetStat("IP")
      ip += team.roster[n].GetStat("IP")
    }
    return avg / ip
  }

  return 0.0
}
