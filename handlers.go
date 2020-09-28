package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// render chooses the way of processing of request, checking it's header
func renderPage(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data)
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}

// ShowLoggingForm renders logging form
func ShowLoggingForm(c *gin.Context) {
	renderPage(c, gin.H{}, "loggingForm.html")
}

// ShowIndexPage renders main page of the project
func ShowIndexPage(c *gin.Context) {
	var Articles []article
	err := DbArticles.Find(nil).All(&Articles)
	if err != nil {
		println(err.Error())
		return
	}
	renderPage(c, gin.H{"Articles": Articles}, "index.html")
}

// Registration procceses data, which was entered by user to registrate
func Registration(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer c.Request.Body.Close()

	var us = user{
		Login:    c.Request.Form.Get("login"),
		Password: c.Request.Form.Get("password"),
		Admin:    false,
	}

	// Can we use c.WriteString()?
	if us.Login == "" {
		c.Writer.WriteString("Login is empty. Please, try again!\n")
		return
	}

	if us.Password == "" {
		c.Writer.WriteString("Password is empty. Please, try again!\n")
		return
	}

	n, err := DbUsers.Find(bson.M{"login": us.Login}).Count()
	if err != nil {
		println(err.Error())
		return
	}

	if n != 0 {
		c.Writer.WriteString("Login is already used by another user. Please, try again!\n")
		return
	}

	hashedPassword, err := HashPassword(us.Password)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	us.Password = hashedPassword

	DbUsers.Insert(us)

	fmt.Println("Registration was successful.")
	ShowIndexPage(c)

}

// Authorisation procceses data, which was entered by user to authorise
func Authorisation(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer c.Request.Body.Close()

	var us = user{
		Login:    c.Request.Form.Get("login"),
		Password: c.Request.Form.Get("password"),
	}

	if us.Login == "" {
		c.Writer.WriteString("Login is empty. Please, try again!\n")
		return
	}

	if us.Password == "" {
		c.Writer.WriteString("Password is empty. Please, try again!\n")
		return
	}

	result := user{}
	DbUsers.Find(bson.M{"login": us.Login}).One(&result)
	if result.Login == "" {
		c.Writer.WriteString("User with such login was not found. Please, try again!\n")
		return
	}

	if !CheckPasswordHash(us.Password, result.Password) {
		c.Writer.WriteString("Invalid password. Please, try again!\n")
		return
	}

	fmt.Println("Authorisation succesful.")

	// Cookie Token set
	c.SetCookie("authorisation_token", us.Login, 30, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"Login": result.Login})

	c.HTML(http.StatusOK, "index.html", gin.H{"Login": result.Login})
}
