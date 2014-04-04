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

  // Start http thread
  http_recv := make(chan string)
  http_send := make(chan fbtron.Simulation)
  go fbtron.RunWebServer(http_send, http_recv)

  for {
    select {
    case <-http_recv:
      // Got a message from the http server. Collect stats from the simulation
      // threads and give the http server the totals.
      var sim_totals fbtron.Simulation
      for n := range sendchannels {
        sendchannels[n] <- "ping"
        sim := <-recvchannels[n]
        sim_totals.Merge(&sim)
      }
      http_send <- sim_totals

      // TODO: If the message contains a list of users who were drafted, tell
      // the simulation threads.

    case <-time.After(3 * time.Second):
      // Just log something to confirm we're still running
      fmt.Println("...")
    }
  }
}
