package main

import (
  "fbtron/fbtron"

  "fmt"
  "runtime"
  "time"
)

func main() {
  fmt.Println("Starting")

  // Tell Go to use all of the processors
  num_cpus := runtime.NumCPU()
  runtime.GOMAXPROCS(num_cpus)

  sendchannels := make([]chan string, num_cpus)
  recvchannels := make([]chan fbtron.Simulation, num_cpus)
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
