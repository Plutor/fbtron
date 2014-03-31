package fbtron

type Player struct {
  name            string
  stats           map[string]float64
  num_seasons     int
  total_wins      int
}

func (p *Player) SetStat(name string, value float64) {
  if p.stats == nil {
    p.stats = make(map[string]float64)
  }
  p.stats[name] = value
}

func (p *Player) GetStat(name string) float64 {
  return p.stats[name]
}

func (p *Player) WinsPerDraft() float64 {
  if p.total_wins == 0 {
    return 0.0
  }
  return float64(p.total_wins) / float64(p.num_seasons)
}

func (p *Player) ResetWins() {
  p.total_wins = 0
  p.num_seasons = 0
}
