package fbtron

const (
  STAT_SUMMED          = 1 << iota
  STAT_IP_WEIGHTED_AVG = 1 << iota
  STAT_AB_WEIGHTED_AVG = 1 << iota

  STAT_ASC             = 1 << iota    // higher is better
  STAT_DESC            = 1 << iota    // lower is better
)

var stat_types = map[string]int {
  "R":    STAT_SUMMED | STAT_ASC,
  "HR":   STAT_SUMMED | STAT_ASC,
  "RBI":  STAT_SUMMED | STAT_ASC,
  "SB":   STAT_SUMMED | STAT_ASC,
  "BA":   STAT_AB_WEIGHTED_AVG | STAT_ASC,

  "W":    STAT_SUMMED | STAT_ASC,
  "S":    STAT_SUMMED | STAT_ASC,
  "K":    STAT_SUMMED | STAT_ASC,
  "ERA":  STAT_IP_WEIGHTED_AVG | STAT_DESC,
  "WHIP": STAT_IP_WEIGHTED_AVG | STAT_DESC,
}

// TODO: Should we return a closure instead?
func GetStatType(name string) int {
  t, ok := stat_types[name]
  if ok {
    return t
  }
  return -1
}
