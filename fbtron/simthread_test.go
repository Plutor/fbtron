package fbtron

import (
  "testing"
)

func TestSimthreadUnimplemented(t *testing.T) {
  t.Log("This function will be removed when tests are implemented")
}

// TODO: Test starting a simulation thread, ask it for data, make sure data is
// good.

func TestRunSeason(t *testing.T) {
  var sim Simulation

  sim.RunSeason()
  if v := sim.Num_seasons; v != 1 {
    t.Errorf("Error running season: expected 1, got %d", v)
  }
}
