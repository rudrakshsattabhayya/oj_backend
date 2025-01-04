package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rudrakshsattabhayya/oj_backend_go/config"
	"github.com/rudrakshsattabhayya/oj_backend_go/helpers"
	"github.com/rudrakshsattabhayya/oj_backend_go/middlewares"
	models "github.com/rudrakshsattabhayya/oj_backend_go/models/session"
	AuthApis "github.com/rudrakshsattabhayya/oj_backend_go/services/auth/apis"
	OjApis "github.com/rudrakshsattabhayya/oj_backend_go/services/oj/apis"
)

func main() {
	config.ConnectDB()
	config.InitSupabaseClient()

	// models.MigrateDB()

	// r := routes.SetupRouter()

	r := gin.Default()

	//If scope of apis.json have admin, then check for admin. Add it to middleware

	r.GET("/dont-sleep", func(c *gin.Context) {
		err := OjApis.InitData()

		if err != nil {
			log.Print("Error in DontSleep:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "ok"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.POST("/Apis/login", func(c *gin.Context) {
		var params AuthApis.LoginParams
		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := AuthApis.Login(params)

		if err != nil {
			log.Print("Error", err.Error())
			c.JSON(400, err.Error())
		} else {
			c.JSON(200, res)
		}
	})

	r.POST("/Apis/login-with-password", func(c *gin.Context) {
		var params AuthApis.LoginWithPasswordParams
		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := AuthApis.LoginWithPassword(params)

		if err != nil {
			log.Print("Error ", err.Error())
			c.JSON(400, err.Error())
		} else {
			c.JSON(200, res)
		}
	})

	r.POST("/Apis/change-the-password", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params AuthApis.ChangeThePasswordParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := AuthApis.ChangeThePassword(params)

		if err != nil {
			log.Print("Error in ChangeThePassword:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/change-username", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params AuthApis.ChangeUsernameParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := AuthApis.ChangeUsername(params)

		if err != nil {
			log.Print("Error in ChangeUsername:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/auth-route", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params AuthApis.AuthenticateRouteParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := AuthApis.AuthenticateRoute(params)

		if err != nil {
			log.Print("Error in Auth Route:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/auth-route-admin", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params AuthApis.AuthenticateRouteAdminParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := AuthApis.AuthenticateRouteAdmin(params)

		if err != nil {
			log.Print("Error in Auth Route:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/create-problem", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params OjApis.CreateProblemParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := helpers.PopulateDynamicStructWithFiles(c, &params); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.CreateProblem(params)

		if err != nil {
			log.Print("Error in Create Problem:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/create_tags", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params OjApis.CreateTagsParams

		_, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.CreateTags(params)

		if err != nil {
			log.Print("Error in Create Problem:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.GET("/Apis/list-problems", func(c *gin.Context) {
		var params OjApis.ListProblemsParams

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.ListProblems(params)

		if err != nil {
			log.Print("Error in List Problems:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/show-problem-solution", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params OjApis.ShowProblemSolutionParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.ShowProblemSolution(params)

		if err != nil {
			log.Print("Error in ShowProblemSolution:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/list-submissions", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params OjApis.ListSubmissionsParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.ListSubmissions(params)

		if err != nil {
			log.Print("Error in ListSubmissions:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.GET("/Apis/list-tags", func(c *gin.Context) {
		var params OjApis.ListTagsParams

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.ListTags(params)

		if err != nil {
			log.Print("Error in ListTags:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/show-problem", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params OjApis.ShowProblemParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.ShowProblem(params)

		if err != nil {
			log.Print("Error in ShowProblem:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.GET("/Apis/get-leaderboard", func(c *gin.Context) {
		var params OjApis.GetLeaderboardParams

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.GetLeaderboard(params)

		if err != nil {
			log.Print("Error in GetLeaderboard:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.POST("/Apis/submit-problem", middlewares.AuthMiddleware, func(c *gin.Context) {
		var params OjApis.SubmitProblemParams

		currSession, exists := c.Get("session")

		if !exists {
			log.Print("Error", "Session does not exist")
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		params.PerformedByID = currSession.(models.Session).PerformedByID

		if err := helpers.PopulateDynamicStructWithFiles(c, &params); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := c.Bind(&params); err != nil {
			log.Print("Error", "Failed to bind parameters:", err.Error())
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		res, err := OjApis.SubmitProblem(params)

		if err != nil {
			log.Print("Error in SubmitProblem:", err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})

	r.Run(":" + os.Getenv("SERVER_PORT"))
}
