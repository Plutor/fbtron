package fbtron

import "fmt"

const (
  STAT_SUMMED          = 1 << iota
  STAT_IP_WEIGHTED_AVG = 1 << iota
  STAT_AB_WEIGHTED_AVG = 1 << iota

  STAT_ASC             = 1 << iota    // higher is better
  STAT_DESC            = 1 << iota    // lower is better
)

var stat_types = map[string]int {
  "B_R":    STAT_SUMMED | STAT_ASC,
  "B_HR":   STAT_SUMMED | STAT_ASC,
  "B_RBI":  STAT_SUMMED | STAT_ASC,
  "B_SB":   STAT_SUMMED | STAT_ASC,
  "B_AVG":  STAT_AB_WEIGHTED_AVG | STAT_ASC,

  "P_W":    STAT_SUMMED | STAT_ASC,
  "P_SV":   STAT_SUMMED | STAT_ASC,
  "P_K":    STAT_SUMMED | STAT_ASC,
  "P_ERA":  STAT_IP_WEIGHTED_AVG | STAT_DESC,
  "P_WHIP": STAT_IP_WEIGHTED_AVG | STAT_DESC,
}

// TODO: Should we return a closure instead?
func GetStatNameAndType(pos, name string) (string, int) {
  fullname := fmt.Sprintf("%s_%s", pos, name)

  return fullname, GetStatType(fullname)
}

func GetStatType(name string) (int) {
  t, ok := stat_types[name]
  if ok {
    return t
  }
  return -1
}
