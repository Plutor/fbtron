package main

import (
  "fbtron/fbtron"

  "fmt"
  "runtime"
  "time"
)

func main() {
  fmt.Println("Starting")
  start_time := time.Now()

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
    time.Sleep(time.Second * 1)

    for n := range sendchannels {
      sendchannels[n] <- "ping"
      response := <-recvchannels[n]
      fmt.Printf("Got response from simulation thread %d: ran %d seasons\n",
                  n, response.Num_seasons)
    }

    if time.Since(start_time) > time.Second * 10 {
      fmt.Println("Shutting down all threads")
      for n := range sendchannels {
        sendchannels[n] <- "quitquitquit"
      }
      break
    }
  }
}
