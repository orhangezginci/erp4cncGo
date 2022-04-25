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

type CreateManufacturerInput struct {
Name string `json:"Name"`
Info       string `json:"info"`
}

type Manufacturer struct {
	Name string
	Info string
	}
func GetManus(c *gin.Context) {
	var MFs []models.Manufacturer
	models.DB.Where("Delete = ?", false).Order("id").Find(&MFs,)
	c.JSON(http.StatusOK, gin.H{"data": MFs})
}

func CreateManu(c *gin.Context) {
	var manu Manufacturer
	println("CreatingManu()")
	jsonbody, _ := ioutil.ReadAll(c.Request.Body)
    //println(string(jsonbody))
	json.Unmarshal([]byte(jsonbody), &manu)
	println("Name:"+manu.Name)
	println("Info:"+manu.Info)


	manudb := models.Manufacturer{Name: manu.Name, Info: manu.Info}
	db := models.DB//c.MustGet("db").(*gorm.DB)
	db.Create(&manudb)
	c.JSON(http.StatusOK, gin.H{"data": manudb})
}

func DeleteManu(c *gin.Context){
	db := models.DB//c.MustGet("db").(*gorm.DB)
	id := (c.Param("manuId"))
  upd_manu := models.Manufacturer{}
	db.Find(&upd_manu,id)
	upd_manu.Delete = true;
	db.Save(&upd_manu)

}
