package main

// InitialiseRouters keeps all routers used in the application
func InitialiseRouters() {
	router.GET("/", ShowIndexPage)

	router.GET("/loggingForm", ShowLoggingForm)

	router.POST("/indexPage", ShowIndexPage)

	router.POST("/registration", Registration)

	router.POST("/authorisation", Authorisation)

	AuthorisedUser := router.Group("/authorisation")
	{
		AuthorisedUser.Use(CheckRegisteredUser)
	}

}
