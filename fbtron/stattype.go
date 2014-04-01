package fbtron

const (
  STAT_SUMMED = iota
  STAT_IP_WEIGHTED_AVG = iota
  STAT_AB_WEIGHTED_AVG = iota
)

var stat_types = map[string]int {
  "R": STAT_SUMMED,
  "HR": STAT_SUMMED,
  "RBI": STAT_SUMMED,
  "SB": STAT_SUMMED,
  "BA": STAT_AB_WEIGHTED_AVG,

  "W": STAT_SUMMED,
  "S": STAT_SUMMED,
  "K": STAT_SUMMED,
  "ERA": STAT_IP_WEIGHTED_AVG,
  "WHIP": STAT_IP_WEIGHTED_AVG,
}

// TODO: Should we return a closure instead?
func GetStatType(name string) int {
  t, ok := stat_types[name]
  if ok {
    return t
  }
  return -1
}
