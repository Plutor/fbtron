package fbtron

import "testing"

func TestGetStatType(t *testing.T) {
  for stat, exp_type := range stat_types {
    if v := GetStatType(stat); v != exp_type {
      t.Errorf("GetStatType: expected type %d for stat '%s', got %d",
               exp_type, stat, v)
    }
  }

  // Unknown type
  if v := GetStatType("ZOMGBBQ"); v != -1 {
    t.Errorf("GetStatType: expected -1 for unknown stat, got %d", v)
  }
}