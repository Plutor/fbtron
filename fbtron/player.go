package fbtron

import (
  "encoding/csv"
  "fmt"
  "io"
  "math/rand"
  "os"
  "strconv"
  "strings"
)

type Player struct {
  ID              string
  Firstname       string
  Lastname        string
  Positions       []string
  stats           map[string]float64
  Num_seasons     int
  Total_wins      int
}

func (p *Player) SetStat(name string, value float64) {
  if p.stats == nil {
    p.stats = make(map[string]float64, len(stat_types))
  }
  p.stats[name] = value
}

func (p *Player) GetStat(name string) float64 {
  return p.stats[name]
}

func (p *Player) GetName() string {
  if p.Lastname == "" {
    return p.Firstname
  }
  return fmt.Sprintf("%s %s", p.Firstname, p.Lastname)
}

func (p *Player) WinsPerDraft() float64 {
  if p.Total_wins == 0 {
    return 0.0
  }
  return float64(p.Total_wins) / float64(p.Num_seasons)
}

func (p *Player) ResetWins() {
  p.Total_wins = 0
  p.Num_seasons = 0
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
    val, err := strconv.ParseFloat(record[n], 64)

    switch header[n] {
    case "steamerid":
      p.ID = record[n]
    case "firstname":
      p.Firstname = record[n]
    case "lastname":
      p.Lastname = record[n]
    case "position":
      position := strings.ToUpper(record[n])
      if position == "LF" || position == "CF" || position == "RF" {
        position = "OF"
      }
      p.Positions = []string { position }
    case "start_percent":
      // TODO: I don't like having to do this, but since the 2013 steamer files
      // don't have eligible positions, we've gotta fake it. The overlap is so
      // that some players qualify as both SP and RP.
      p.Positions = []string { }
      if val >= 0.25 {
        p.Positions = append(p.Positions, "SP")
      }
      if val <= 0.75 {
        p.Positions = append(p.Positions, "RP")
      }
    case "IP":
      if err == nil {
        p.SetStat("P_IP", val)
      }
    case "AB":
      if err == nil {
        p.SetStat("B_AB", val)
      }
    default:
      stat_name, stat_type := GetStatNameAndType(default_position, header[n])
      if stat_type != -1 {
        if err == nil {
          p.SetStat(stat_name, val)
        }
      }
    }
  }

  // Add default position if it's non-blank and isn't already a position
  if default_position != "" {
    for n := range p.Positions {
      if p.Positions[n] == default_position {
        return p
      }
    }
    p.Positions = append(p.Positions, default_position)
  }

  // If the player doesn't have an id, generate one randomly (there is no
  // guarantee these won't collide, but chances are 1/2^64).
  if p.ID == "" {
    p.ID = fmt.Sprintf("%x", rand.Int63())
  }

  return p
}
