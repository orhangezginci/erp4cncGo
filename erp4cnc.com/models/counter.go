package models

import "time"

type Counter struct {
   ID    uint64 `json:"id" gorm:"primaryKey"`

  Machine_id  uint `json:"machineId"`
  Auftrags_id uint `json:"auftragsId"`
  Mitarbeiter_id uint `json:"mitarbeiterId"`
  Count uint `json:"count"`
  CreatedAt time.Time `json:"createdAt" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
  Diff int `json:"diff"`
}
