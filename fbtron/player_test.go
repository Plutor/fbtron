package fbtron

import (
  "testing"
)

func TestSetStat(t *testing.T) {
  p := new(Player)
  p.SetStat("test", 1.0)
  if v:= p.GetStat("test"); v != 1.0 {
    t.Errorf("Failure to set/get a stat: expected 1.0, got %f", v)
  }
}

func TestWinsPerDraft(t *testing.T) {
  p := new(Player)
  p.Num_seasons = 10
  if v := p.WinsPerDraft(); v != 0.0 {
    t.Errorf("Failure to get WinsPerDraft: expected 0.0, got %f", v)
  }

  p.Total_wins = 15
  if v := p.WinsPerDraft(); v != 1.5 {
    t.Errorf("Failure to get WinsPerDraft: expected 1.5, got %f", v)
  }

  p.ResetWins()
  if v := p.WinsPerDraft(); v != 0.0 {
    t.Errorf("Failure to get WinsPerDraft: expected 0.0, got %f", v)
  }
}

func TestBuildPlayersFromCsv(t *testing.T) {
  var players []*Player

  // Pass empty file, expect empty array
  players = BuildPlayersFromCsv("testdata/players_csv_empty.csv", "")
  if len(players) > 0 {
    t.Errorf("BuildPlayersFromCsv: expected empty array, got %s", players)
  }

  // Pass file with just a header, expect empty array
  players = BuildPlayersFromCsv("testdata/players_csv_headeronly.csv", "")
  if len(players) > 0 {
    t.Errorf("BuildPlayersFromCsv: expected empty array, got %s", players)
  }

  // Pass broken csv, expect empty array
  players = BuildPlayersFromCsv("testdata/players_csv_broken.csv", "")
  if len(players) > 0 {
    t.Errorf("BuildPlayersFromCsv: expected empty array, got %s", players)
  }

  // Pass csv with a header and a record, expect one-player array back
  players = BuildPlayersFromCsv("testdata/players_csv_ok.csv", "B")
  if len(players) != 5 {
    t.Errorf("BuildPlayersFromCsv: expected 5 players, got %d:\n%s",
             len(players), players)
  } else {
    if v := players[0].GetName(); v != "Foo Bar" {
      t.Errorf("BuildPlayerFromCsvRecord: expected name Foo Bar, got '%s'", v)
    }
    if v := players[0].GetStat("B_R"); v != 100 {
      t.Errorf("BuildPlayerFromCsv: expected R=100, got %f", v)
    }
    if v := players[0].GetStat("B_RBI"); v != 200 {
      t.Errorf("BuildPlayerFromCsv: expected RBI=200, got %f", v)
    }
    if v := players[0].Positions; len(v) != 2 || v[0] != "2B" || v[1] != "B" {
      t.Errorf("BuildPlayerFromCsv: expected position 2B,B, got %s", v)
    }
  }
}

func TestBuildPlayerFromCsvRecord(t *testing.T) {
  var player *Player

  // Pass empty arrays, expect nil
  player = BuildPlayerFromCsvRecord([]string {}, []string {}, "")
  if (player != nil) {
    t.Errorf("BuildPlayerFromCsvRecord: expected nil, got %s",
             player)
  }

  // Pass header array but empty data, expect nil
  player = BuildPlayerFromCsvRecord(
      []string {"firstname", "lastname", "R", "RBI"},
      []string {},
      "")
  if (player != nil) {
    t.Errorf("BuildPlayerFromCsvRecord: expected nil, got %s",
             player)
  }

  // Pass data but empty header, expect nil
  player = BuildPlayerFromCsvRecord(
      []string {},
      []string {"Foo", "Bar", "100", "200"},
      "")
  if (player != nil) {
    t.Errorf("BuildPlayerFromCsvRecord: expected nil, got %s",
             player)
  }

  // Pass data and header, expect Player
  player = BuildPlayerFromCsvRecord(
      []string {"firstname", "lastname", "R", "RBI"},
      []string {"Foo", "Bar", "100", "200"},
      "B")
  if (player == nil) {
    t.Errorf("BuildPlayerFromCsvRecord: expected Player, got nil")
  }
  if v := player.GetName(); v != "Foo Bar" {
    t.Errorf("BuildPlayerFromCsvRecord: expected name Foo Bar, got '%s'", v)
  }
  if v := player.GetStat("B_R"); v != 100 {
    t.Errorf("BuildPlayerFromCsvRecord: expected R=100, got %f", v)
  }
  if v := player.GetStat("B_RBI"); v != 200 {
    t.Errorf("BuildPlayerFromCsvRecord: expected RBI=200, got %f", v)
  }
  if v := player.Positions; len(v) != 1 || v[0] != "B" {
    t.Errorf("BuildPlayerFromCsvRecord: expected position B, got %s", v)
  }

  // TODO: Player with same position as default
  // TODO: Empty string as default
}
