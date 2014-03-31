package fbtron

import (
  "testing"
)

func TestPlayerUnimplemented(t *testing.T) {
  t.Log("This function will be removed when tests are implemented")
}

func TestSetStat(t *testing.T) {
  p := new(Player)
  p.SetStat("test", 1.0)
  if v:= p.GetStat("test"); v != 1.0 {
    t.Errorf("Failure to set/get a stat: expected 1.0, got %f", v)
  }
}
