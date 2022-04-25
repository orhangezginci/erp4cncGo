// controllers/manufacturer.go
package controllers
import (
// "bookc/models"
"gezginci.com/erp4cnc/models"
"net/http"
_"fmt"
"io/ioutil"
"encoding/json"
"github.com/gin-gonic/gin"
_"github.com/jinzhu/gorm"
)

type Machine struct {
	Name string `json:"Name"`
	Info string `json:"Info"`
	Rev1 bool `json:"Rev1"`
	Rev2 bool `json:"Rev2"`
	Rev3 bool `json:"Rev3"`
	Frsspindel bool `json:"Frsspindel"`
	Manufacturer uint `json:"Manufacturer"`
	Type uint  `json:"Type"`
}

func CreateMachine(c *gin.Context) {
	var machine Machine
	println("CreateMachine()")
	jsonbody, _ := ioutil.ReadAll(c.Request.Body)
  println(string(jsonbody))
///////////////////////////////////////////////////
	jsonMap := make(map[string]interface{})
  err := json.Unmarshal([]byte(jsonbody), &jsonMap)
	println(jsonMap)
	println(err)
///////////////////////////////////////////////////

	json.Unmarshal([]byte(jsonbody), &machine)
	println("Name:"+machine.Name)
	println()
  println(machine.Rev1,machine.Rev2,machine.Rev3)
	machinedb := models.Machine{Name: machine.Name,Info:machine.Info, Type: machine.Type,Manufacturer: machine.Manufacturer,
									Rev1:machine.Rev1,Rev2:machine.Rev2,Rev3:machine.Rev3,Frsspindel:machine.Frsspindel}
	db := models.DB//c.MustGet("db").(*gorm.DB)
	db.Create(&machinedb)
	c.JSON(http.StatusOK, gin.H{"data": machine})
}

func GetMachines(c *gin.Context) {
	println("GetMachines")
	var machines []models.Machine
	models.DB.Where("Deleted = ?", false).Order("id").Find(&machines)
	c.JSON(http.StatusOK, gin.H{"data": machines})
}

/*
type Machine struct {
  ID   uint   `json:"id" gorm:"primary_key"`
  Name string `json:"name"`
  Type uint `json:"type"`
  Manufacturer uint `json:"manufacturer"`
  Rev1 bool `json:"rev1"`
  Rev2 bool `json:"rev2"`
  Rev3 bool `json:"rev3"`
  Frsspindel bool `json:"frsspindel"`
}*/
