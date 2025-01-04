package controllers

// import (
// 	"reflect"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rudrakshsattabhayya/oj_backend_go/services/auth/apis"
// )

// var paramsMap = map[string]reflect.Type{
// 	"SignUpParams": reflect.TypeOf(apis.SignUpParams{}),
// }

// var apisMap = map[string]interface{}{
// 	"SignUp": apis.SignUp,
// }


// func AuthController(c *gin.Context, endpoint string) {
// 	api := apisMap[endpoint]
// 	params := reflect.New(reflect.TypeOf(paramsMap[endpoint + "Params"])).Interface()

// 	err := c.Bind(&params)
// 	if err != nil {
// 		log.Print("Error", err.Error())
// 		c.JSON(c.Writer.Status(), err.Error())
// 		return
// 	}

// 	ser, err := api.Call(id)

// }
