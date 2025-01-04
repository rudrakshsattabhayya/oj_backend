package routes

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rudrakshsattabhayya/oj_backend_go/middlewares"
// )

// type RouteInfo struct {
// 	Method     string `json:"method"`
// 	AccessType string `json:"access_type"`
// }

// type Routes map[string]RouteInfo

// func LoadRoutes(folder string) (map[string]Routes, error) {
// 	files, err := os.ReadDir(folder)
// 	if err != nil {
// 		return nil, err
// 	}

// 	routes := make(map[string]Routes)
// 	for _, file := range files {
// 		if filepath.Ext(file.Name()) == ".json" {
// 			filePath := filepath.Join(folder, file.Name())
// 			content, err := os.ReadFile(filePath)
// 			if err != nil {
// 				return nil, err
// 			}

// 			var route Routes
// 			err = json.Unmarshal(content, &route)
// 			if err != nil {
// 				return nil, err
// 			}

// 			serviceName := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
// 			routes[serviceName] = route
// 		}
// 	}

// 	return routes, nil
// }

// func getPascalCasedString(s string) string {
// 	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(s, "_", " ")), " ", "")
// }

// func setupRouteHandlers(router *gin.Engine, serviceName string, routes Routes) {
// 	for endpoint, info := range routes {
// 		fullPath := fmt.Sprintf("/%s/%s", serviceName, endpoint)

// 		pascalCasedEndpoint := getPascalCasedString(endpoint)

// 		switch info.Method {
// 		case "get":
// 			if info.AccessType == "private" {
// 				router.GET(fullPath, middlewares.AuthMiddleware, handlerFunc)
// 			} else {
// 				router.GET(fullPath, handlerFunc)
// 			}
// 		case "post":
// 			if info.AccessType == "private" {
// 				router.POST(fullPath, middlewares.AuthMiddleware, handlerFunc)
// 			} else {
// 				router.POST(fullPath, handlerFunc)
// 			}
// 		}
// 	}
// }

// func SetupRouter() *gin.Engine {
// 	router := gin.Default()

// 	// Health check route
// 	router.GET("/health", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"message": "ok"})
// 	})

// 	// Load routes from JSON files in the same folder
// 	routes, err := LoadRoutes("./routes")
// 	if err != nil {
// 		fmt.Printf("Error loading routes: %v\n", err)
// 		return router
// 	}

// 	fmt.Println(routes)

// 	// Dynamically add routes
// 	for serviceName, serviceRoutes := range routes {
// 		setupRouteHandlers(router, serviceName, serviceRoutes)
// 	}

// 	return router
// }
