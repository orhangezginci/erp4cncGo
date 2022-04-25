package main

import "github.com/gin-gonic/gin"
import "gezginci.com/erp4cnc/models" // new
import "gezginci.com/erp4cnc/controllers" // new
import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"net/http"
	"os"
	"time"
	"github.com/go-redis/redis/v7"
  	"github.com/twinj/uuid"
	  "fmt"
	  "strconv"
	  "strings"
  )

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
  }
  type AccessDetails struct {
    AccessUuid string
    UserId   uint64
}
  type User struct {
	ID uint64       `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone string 	`json:"phone"`
  }
  var user = User{
	ID:            1,
	Username: "username",
	Password: "password",
	Phone: "49123454322", //this is a random number
  }
  var (
	 client *redis.Client
	 router = gin.Default()
)
func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	fmt.Println("DSN::"+dsn)
	if len(dsn) == 0 {
	   dsn = "redis1:6379"
	}
	client = redis.NewClient(&redis.Options{
	   Addr: dsn, //redis port
	})
	ret, err := client.Ping().Result()
	if err != nil {
	   panic(err)
	} else {
		fmt.Println(ret)
	}
  }

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3001"}


	models.ConnectDatabase() // new
	models.DB.AutoMigrate(&models.Manufacturer{})
	models.DB.AutoMigrate(&models.Machine{})

	r.GET("/test",test)
	r.POST("/answer",answer)
	r.POST("/login", Login)
	r.GET("/ping",TokenAuthMiddleware_(), pong)
	r.POST("/API/ManuFacture",controllers.CreateManu)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
func test(c *gin.Context){
	fmt.Println(">>>>>>>>>>>>>>>>>>>TEST!!!!!!!!!!")
	/*c.HTML(200, "dashboard", gin.H{
        "title":    "Dashboard",
        })*/
}
func pong(c *gin.Context){
	fmt.Println("PONG!!!!!!!!")
	/*c.HTML(200, "dashboard", gin.H{
        "title":    "Dashboard",
        })*/
}
func cookie(c *gin.Context) {
	cookie,err := c.Request.Cookie("jwt");
	if(cookie!= nil){
	   fmt.Println("COOKIE NOT NULL")

	   fmt.Println("func cookie(c *gin.Context):::::"+cookie.Value+":::::")
	   }
	if(err!= nil){fmt.Println("error:"+err.Error())}
}
func getCookieValue(c *gin.Context) (string) {
	cookie,err := c.Request.Cookie("jwt");
	for _, cookie := range c.Request.Cookies() {
		fmt.Println("Found a cookie named:", cookie.Name)
		fmt.Println("Value of Cookie:",cookie.Value)
	  }

	if(cookie!= nil){
	fmt.Println("COOKIE NOT NULL")
	fmt.Println("func getCookieValue(c *gin.Context):::::"+cookie.Value+":::::")
	return cookie.Value
	}
	if(err!= nil){
		fmt.Println("error:"+err.Error())
		return ""
	}
	return ""
}
func answer(c *gin.Context) {

	expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
    http.SetCookie(c.Writer, &cookie)
	//http.SetCookie("testhttpcookie", "testvalue", 1000, "/", "127.0.0.1", false, false)

    /*c.HTML(200, "dashboard", gin.H{
        "title":    "Dashboard",
        })*/
}
func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
	   c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	   return
	}
	//compare the user from the request, with the one we defined:
	fmt.Println("Username:"+u.Username)
	fmt.Println("Password:"+u.Password)
	if user.Username != u.Username || user.Password != u.Password {
		 fmt.Println("Please provide valid login details")
	   c.JSON(http.StatusUnauthorized, "Please provide valid login details")
	   return
	}
	answer(c)
	fmt.Println("CreateToken:%d",user.ID)
	ts, err := CreateToken(user.ID)
   if err != nil {
	fmt.Println("CreateToken-2:%d",user.ID)
   c.JSON(http.StatusUnprocessableEntity, err.Error())
	 return
	 exp := time.Now().Add(365 * 24 * time.Hour)
	 cookie := http.Cookie{Name: "jwt", Value: ts.AccessToken,Expires: exp}
	 http.SetCookie(c.Writer, &cookie)
	}

	fmt.Println("CreateAuth:%d",user.ID)

   saveErr := CreateAuth(user.ID, ts)
	if saveErr != nil {
		fmt.Println("StatusUnprocessableEntity:%d",user.ID)

	   c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
	   "refresh_token": ts.RefreshToken,
	}
	if(tokens!= nil){
		fmt.Printf("TOKEN::::))))"+ts.AccessToken)
		//c.SetCookie("jwt",ts.AccessToken, 1000, "/", "", true, false)
		exp := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "jwt", Value: ts.AccessToken,Expires: exp}
    http.SetCookie(c.Writer, &cookie)
	}
  }
  func CreateToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
	   return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
	   return nil, err
	}
	return td, nil
  }
  func CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
	 return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
	 return errRefresh
	}
	return nil
   }
   func ExtractToken(c *gin.Context) string {
	fmt.Println("ExtractToken")
 	//bearToken := r.Header.Get("Authorization")
 	bearToken := getCookieValue(c)
	fmt.Println("bearToken>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>::"+bearToken+"::<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	// bearToken_ := strings.Replace(bearToken,"#:#","*",0)
	// fmt.Println("bearToken>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>__"+bearToken_+"__<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

 //normally Authorization the_token_xxx
 bearToken="Token "+bearToken
 strArr := strings.Split(bearToken, " ")

 if len(strArr) == 2 {
	return strArr[1]
 }

 return ""
}
func VerifyTokenGin(r *http.Request,c *gin.Context) (*jwt.Token, error) {
	fmt.Println("VerifyToken")
	tokenString := ExtractToken(c)
	//fmt.Println("tokenString>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"+tokenString+"<<<<<<<<<<<<")
	//check claims and checkDB
	fmt.Println("VerifyToken-2")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, "Authorization failed, No Token in Header!")
		c.Abort()
		var jwtToken *jwt.Token
		return jwtToken, nil
	 }


	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
 fmt.Println("jwt.Parse################################",tokenString)

 claims, _ := token.Claims.(jwt.MapClaims)
		access_uuid := claims["access_uuid"].(string)
		fmt.Println("access_uuid",access_uuid)
		//check access_uuid in Db
		if(!isSessionInDb(access_uuid)){
			c.JSON(http.StatusUnauthorized, "Authorization failed, No Valid Session!")
		c.Abort()

			//fmt.Println("Authorization failed, No Valid Session")
			//c.JSON(http.StatusUnauthorized,"No Valid Session")


			//return nil, fmt.Errorf("Authorization failed No Valid Session Dick)
		}

	   //Make sure that the token method conform to "SigningMethodHMAC"
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		fmt.Println("VerifyToken-4")
		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	   }
	   fmt.Println("VerifyToken-5")
	   fmt.Println("VerifyToken-5-1")
	   fmt.Println("VerifyToken-5-2")
	   fmt.Println("VerifyToken-5-3")
	   fmt.Println("VerifyToken-5-4")
	   lolo:= []byte("jdnfksdmfksd")
	   fmt.Println("VerifyToken-5-5")
	   fmt.Println(lolo)

	   //return []byte(os.Getenv("ACCESS_SECRET")), nil
	   return []byte("jdnfksdmfksd"), nil

	})
	if err != nil {
		fmt.Println("JWT SIGNING METHOD",token.Header["alg"])
		  fmt.Println("VerifyToken-3:::::"+err.Error())

	   return nil, err
	}
	return token, nil
  }
  func isSessionInDb(uuid string)(bool){
	_, err := client.Get(uuid).Result()
	if (err != nil) {
		fmt.Println("Error fetching from DB",err.Error())
		return false
	}else {
	return true}}

	func TokenValid(r *http.Request,c *gin.Context) error {
		fmt.Println("TokenValid")
		token, err := VerifyTokenGin(r,c)
		if(err != nil){
		fmt.Println("VerifyToken:",err.Error())
		}
		fmt.Println("TokenValid-2")

		if err != nil {
			fmt.Println("TokenValid-3")
		   return err
		}
		fmt.Println("TokenValid-4")
		if(token!=nil){
			fmt.Println("TOKEN IS NOT NULL")
			if (token.Valid){
				fmt.Println("TOKEN IS VALID")
			}

			} else{
				fmt.Println("INVALID TOKEN")
				return err
			}
		if claims, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
			fmt.Println("TokenValid-5 claims",claims["access_uuid"])
			return nil
		}
		fmt.Println("TokenValid-6")
		return nil
	  }
	  func FetchAuth(authD *AccessDetails) (uint64, error) {
		fmt.Println("FetchAuth")

		userid, err := client.Get(authD.AccessUuid).Result()
		if err != nil {
			fmt.Println("err>>>>>>>>><<<<<<>>>>"+err.Error())
		   return 0, err
		}
		userID, _ := strconv.ParseUint(userid, 10, 64)
		return userID, nil
	  }
	  func TokenAuthMiddleware__() gin.HandlerFunc {
		//c:=*gin.Context()
	  return func(c *gin.Context) {
		  fmt.Println("TokenAuthMiddleware-----------")

		  err:= TokenValid(c.Request,c)

		  if err != nil {
			  //fmt.Println("ERROR:>>>>>"+err.Error())
			  c.Abort()
		  } else{
			  //fmt.Println("NOERROR:>>>>>")

			 c.Next()
		  }
	  /*c.Next()*/}
	  }
	  func TokenAuthMiddleware_() gin.HandlerFunc {
		return func(c *gin.Context) {
			fmt.Println("TokenAuthMiddleware")
		   err := TokenValid(c.Request,c)
		   getCookieValue(c)

		   if err != nil {
			  c.JSON(http.StatusUnauthorized, err.Error())
			  c.Abort()
			  return
		   }
		   //c.Next()
		}
	  }
