package main

import (
  "fbtron/fbtron"

  "testing"
  "time"
)

func TestRunSimulation(t *testing.T) {
  sendchannel := make(chan string, 1)
  recvchannel := make(chan fbtron.Simulation, 1)
  go fbtron.RunSimulation(sendchannel, recvchannel)

  // TODO: Sleeping is frowned upon, do it a different way.
  time.Sleep(250 * time.Millisecond)

  sendchannel <- "ping"
  data := <-recvchannel
  if data.Num_seasons == 0 {
    t.Errorf("RunSimulation: Expected Simulation with >0 seasons, got %d",
             data.Num_seasons)
  }

  sendchannel <- "quitquitquit"
}

func BenchmarkRunSimulation(b *testing.B) {
  // Start a simulation and stop it immediately.
  sendchannel := make(chan string, 1)
  recvchannel := make(chan fbtron.Simulation, 1)
  go fbtron.RunSimulation(sendchannel, recvchannel)
  sendchannel <- "ping"
  sim := <-recvchannel
  sendchannel <- "quitquitquit"

  // Then run N seasons on this real simulation
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    sim.RunSeason()
  }
}
