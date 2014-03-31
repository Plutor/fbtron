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

func TestWinsPerDraft(t *testing.T) {
  p := new(Player)
  p.num_seasons = 10
  if v := p.WinsPerDraft(); v != 0.0 {
    t.Errorf("Failure to get WinsPerDraft: expected 0.0, got %f", v)
  }

  p.total_wins = 15
  if v := p.WinsPerDraft(); v != 1.5 {
    t.Errorf("Failure to get WinsPerDraft: expected 1.5, got %f", v)
  }

  p.ResetWins()
  if v := p.WinsPerDraft(); v != 0.0 {
    t.Errorf("Failure to get WinsPerDraft: expected 0.0, got %f", v)
  }
}
