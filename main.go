package main

import (
  "fbtron/fbtron"

  "flag"
  "fmt"
  "runtime"
  "time"
)

var num_threads = flag.Int(
    "threads", 0, "Number of threads to start (0 = number of cpus)")
var num_cpus = flag.Int(
    "cpus", 0, "Number of CPUs to use (0 = all on system)")

func main() {
  // Tell Go to use all of the processors
  if *num_cpus == 0 {
    *num_cpus = runtime.NumCPU()
  }
  runtime.GOMAXPROCS(*num_cpus)
  if *num_threads == 0 {
    *num_threads = *num_cpus
  }

  fmt.Printf("Starting %d threads on %d cpus\n", *num_threads, *num_cpus)

  sendchannels := make([]chan string, *num_threads)
  recvchannels := make([]chan fbtron.Simulation, *num_threads)
  for n := range sendchannels {
    sendchannels[n] = make(chan string, 1)
    recvchannels[n] = make(chan fbtron.Simulation, 1)
    go fbtron.RunSimulation(sendchannels[n], recvchannels[n])
  }

  // TODO: Start http thread ...

  // TEMP
  for {
    time.Sleep(time.Second * 3)

    for n := range sendchannels {
      sendchannels[n] <- "ping"
      sim := <-recvchannels[n]

      var bestplayer *fbtron.Player
      var bestwpd float64

      for _, player := range sim.All_players {
        if bestplayer == nil || player.WinsPerDraft() > bestwpd {
          bestplayer = player
          bestwpd = player.WinsPerDraft()
        }
      }

      fmt.Printf("t%d: ran %d seasons, best player is %s (%.1f)\n",
                  n, sim.Num_seasons, bestplayer.GetName(), bestwpd)
    }
  }
}
