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

  sendchannels := make([]chan fbtron.UserAction, *num_threads)
  recvchannels := make([]chan fbtron.Simulation, *num_threads)
  for n := range sendchannels {
    sendchannels[n] = make(chan fbtron.UserAction, 1)
    recvchannels[n] = make(chan fbtron.Simulation, 1)
    go fbtron.RunSimulation(sendchannels[n], recvchannels[n])
  }

  // Start http thread
  http_recv := make(chan fbtron.UserAction)
  http_send := make(chan fbtron.Simulation)
  go fbtron.RunWebServer(http_send, http_recv)

  for {
    select {
    case msg := <-http_recv:
      // Got a message from the http server, send it to the simulators, collect
      // and sum the responses, and send it back to the http server.
      var sim_totals fbtron.Simulation
      for n := range sendchannels {
        sendchannels[n] <- msg
        sim := <-recvchannels[n]
        sim_totals.Merge(&sim)
      }
      http_send <- sim_totals
    case <-time.After(30 * time.Second):
      // Just log something to confirm we're still running
      fmt.Println("...")
    }
  }
}
