package models
import "time"

type MachineType struct {
  ID     uint   `json:"id" gorm:"primary_key"`
  Name  string `json:"name"`
  Info string `json:"info"`
  Delete bool `json:"delete"`
  CreatedAt time.Time `json:"createdAt" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`

}
