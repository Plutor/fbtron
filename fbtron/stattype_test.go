package fbtron

import "testing"

func TestGetStatType(t *testing.T) {
  for stat, exp_type := range stat_types {
    v := GetStatType(stat)
    if v != exp_type {
      t.Errorf("GetStatType: expected type %d for stat '%s', got %d",
               exp_type, stat, v)
    }

    // Also every stat type should consist of an accumulator and a direction.
    if v & (STAT_SUMMED | STAT_IP_WEIGHTED_AVG | STAT_AB_WEIGHTED_AVG) == 0 {
      t.Errorf("GetStatType: type %d for stat '%s' doesn't have an accumulator",
               v, stat)
    }
    if v & (STAT_ASC | STAT_DESC) == 0 {
      t.Errorf("GetStatType: type %d for stat '%s' doesn't have a direction",
               v, stat)
    }
  }

  // Unknown type
  if v := GetStatType("ZOMGBBQ"); v != -1 {
    t.Errorf("GetStatType: expected -1 for unknown stat, got %d", v)
  }
}
