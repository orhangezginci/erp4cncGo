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

type MachineType struct {
	Name string
	Info string
	}
func GetTypes(c *gin.Context) {
	var MTs []models.MachineType
	models.DB.Find(&MTs)
	c.JSON(http.StatusOK, gin.H{"data": MTs})
}
func CreateType(c *gin.Context) {
	var manType MachineType



	println("CreatingManu()")
	jsonbody, _ := ioutil.ReadAll(c.Request.Body)
    //println(string(jsonbody))
	json.Unmarshal([]byte(jsonbody), &manType)
	println("Name:"+manType.Name)
	println("Info:"+manType.Info)


	manTypedb := models.MachineType{Name: manType.Name, Info: manType.Info}
	db := models.DB//c.MustGet("db").(*gorm.DB)
	db.Create(&manTypedb)
	c.JSON(http.StatusOK, gin.H{"data": manType})
}
