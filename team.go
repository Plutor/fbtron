package fbtron

type Team struct {
  name    string
  roster  []*Player
  wins    int
}

func (team *Team) AddPlayer(p *Player) {
  nthplayer := len(team.roster) + 1
  if nthplayer > cap(team.roster) {
    newroster := make([]*Player, nthplayer*2)
    copy(team.roster, newroster)
    team.roster = newroster
  }

  team.roster = team.roster[0:nthplayer]
  team.roster[nthplayer-1] = p
}

//func (team *Team) Get