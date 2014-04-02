package fbtron

import (
  "encoding/csv"
  "fmt"
  "io"
  "os"
  "strconv"
)

type Player struct {
  firstname       string
  lastname        string
  positions       []string
  stats           map[string]float64
  num_seasons     int
  total_wins      int
}

func (p *Player) SetStat(name string, value float64) {
  if p.stats == nil {
    p.stats = make(map[string]float64)
  }
  p.stats[name] = value
}

func (p *Player) GetStat(name string) float64 {
  return p.stats[name]
}

func (p *Player) GetName() string {
  if p.lastname == "" {
    return p.firstname
  }
  return fmt.Sprintf("%s %s", p.firstname, p.lastname)
}

func (p *Player) WinsPerDraft() float64 {
  if p.total_wins == 0 {
    return 0.0
  }
  return float64(p.total_wins) / float64(p.num_seasons)
}

func (p *Player) ResetWins() {
  p.total_wins = 0
  p.num_seasons = 0
}

// BuildPlayersFromCsv reads a CSV file and returns an array of player objects,
// one for each row. Assumes the first row is labels. Every column whose label
// is defined in StatType has its value set.
func BuildPlayersFromCsv(filename string, default_position string) []*Player {
  file, err := os.Open(filename)
  if err != nil {
    fmt.Println("Error:", err)
    return []*Player {}
  }
  defer file.Close()
  reader := csv.NewReader(file)

  var header []string
  var players []*Player
  for {
    record, err := reader.Read()
    if err == io.EOF {
      break
    } else if err != nil {
      fmt.Printf("Error reading %s: %s\n", filename, err)
      continue
    }

    if header == nil {
      header = record
    } else {
      newplayer := BuildPlayerFromCsvRecord(header, record, default_position)
      if newplayer != nil {
        if len(players) == cap(players) {
          newplayers := make([]*Player, len(players), (len(players)+1)*2)
          copy(newplayers, players)
          players = newplayers
        }
        players = players[:len(players)+1]
        players[len(players)-1] = newplayer
      }
    }
  }

  // Assume the caller will close the file handle
  return players
}

func BuildPlayerFromCsvRecord(
    header []string, record []string, default_position string) *Player {
  columns := len(header)
  if len(record) < columns {
    columns = len(record)
  }
  if columns == 0 {
    return nil
  }

  p := new(Player)
  for n := 0; n < columns; n ++ {
    switch header[n] {
    case "firstname":
      p.firstname = record[n]
    case "lastname":
      p.lastname = record[n]
    case "position":
      p.positions = []string { record[n] }
    default:
      if GetStatType(header[n]) != -1 {
        val, err := strconv.ParseFloat(record[n], 64)
        if err == nil {
          p.SetStat(header[n], val)
        }
      }
    }
  }

  p.positions = append(p.positions, default_position)

  return p
}
