package fbtron

type Player struct {
  name string
  stats map[string]float64
}

func (p *Player) SetStat(name string, value float64) {
  if p.stats == nil {
    p.stats = make(map[string]float64)
  }
  p.stats[name] = value
}

func (p *Player) GetStat(name string) float64 {
  val, ok := p.stats[name]
  if ok {
    return val
  }
  return -1.0
}
