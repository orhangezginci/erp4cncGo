// models/machine.go

package models

import "time"


type Machine struct {
  ID   uint   `json:"id" gorm:"primary_key"`
  Name string `json:"name"`
  Info string  `json:"info"`
  Type uint `json:"type"`
  Manufacturer uint `json:"manufacturer"`
  Rev1 bool `json:"rev1"`
  Rev2 bool `json:"rev2"`
  Rev3 bool `json:"rev3"`
  Frsspindel bool `json:"frsspindel"`
  Deleted bool `json:"deleted" gorm:"default:false"`
  CreatedAt time.Time `json:"createdAt" gorm:"type:TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}
