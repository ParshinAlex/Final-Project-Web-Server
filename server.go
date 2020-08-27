package main

// Импорты

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"
)

// Структуры для использования

type safeMap struct {
	v   map[string]string
	mux sync.RWMutex
}

var users = safeMap{v: make(map[string]string)}

var mySigningKey = []byte("secret")

func (mapa *safeMap) Update(login string, password string) {
	mapa.mux.Lock()
	defer mapa.mux.Unlock()
	mapa.v[login] = password
}

func (mapa *safeMap) Delete(login string) {
	mapa.mux.Lock()
	defer mapa.mux.Unlock()
	delete(mapa.v, login)
}

type user struct {
	login    string
	password string
}

type product struct {
	id   int    `json:"id"`
	name string `json: "name"`
}

var products = []product{
	product{id: 1, name: "First"},
	product{id: 2, name: "Second"},
}

// Хеширование паролей

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Хендлеры

var registrationFormHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "registrationForm.html")
})

var authorisationFormHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "authorisationForm.html")
})

var startingHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "authorisationForm.html")
})

var notImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

var statusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and working"))
})

var authorisationHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "authorisation.html")

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error!\n", err.Error())
		return
	}
	defer r.Body.Close()

	var us = user{
		login:    r.Form.Get("login"),
		password: r.Form.Get("password"),
	}

	if _, ok := users.v[us.login]; !ok {
		fmt.Fprintf(w, "There is no user with such login. Please, try again!\n")
		return
	}

	if !checkPasswordHash(us.password, users.v[us.login]) {
		fmt.Fprintf(w, "Invalid password. Please, try again!\n")
		return
	}

	//fmt.Fprintf(w, "Authorisation successful. Welcome back, %s!\n", us.login)
	http.ServeFile(w, r, "mainPage.html")
})

var registrationHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "registration.html")

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error!\n", err.Error())
		return
	}
	defer r.Body.Close()

	var us = user{
		login:    r.Form.Get("login"),
		password: r.Form.Get("password"),
	}

	if us.login == "" {
		fmt.Fprintf(w, "Invalid login, please, try again!\n")
		return
	}
	if _, ok := users.v[us.login]; ok {
		fmt.Fprintf(w, "Login is already used by another user. Please, try again!\n")
		return
	}

	hashedPassword, err := hashPassword(us.password)
	if err != nil {
		fmt.Println("Error!\n", err.Error())
		return
	}
	users.Update(us.login, hashedPassword)
	//fmt.Fprintf(w, "Resistered successfully.\n")
	//w.Write([]byte("Registered successfully.\n"))
	authorisationFormHandler(w, r)
})

var productsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error!\n", err.Error())
		return
	}
	defer r.Body.Close()
	tokenString := r.Form.Get("token")
	token, err := jwt.DecodeSegment(tokenString)
	if err != nil {
		fmt.Println("Error!\n", err.Error())
		return
	}
	fmt.Println(token)
	payload, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
})

var addFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var pr product
	vars := mux.Vars(r)
	name := vars["name"]

	for _, p := range products {
		if p.name == name {
			pr = p
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if pr.name != "" {
		payload, _ := json.Marshal(pr)
		w.Write([]byte(payload))
	} else {
		w.Write([]byte("Product Not Found"))
	}

})

// Токены

var getTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	token.Claims = claims
	tokenString, _ := token.SignedString(mySigningKey)
	w.Write([]byte(tokenString))
})

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// Main

func main() {
	adminHashed, err := hashPassword("admin")
	if err != nil {
		return
	}
	users.Update("admin", adminHashed)
	//token := jwt.New(jwt.SigningMethodHS256)

	r := mux.NewRouter()

	r.Handle("/get-token", getTokenHandler)
	r.Handle("/", startingHandler)
	r.Handle("/status", statusHandler)
	r.Handle("/products", jwtMiddleware.Handler(productsHandler))
	r.Handle("/products/{name}/feedback", jwtMiddleware.Handler(addFeedbackHandler))
	r.Handle("/authorisationForm", authorisationFormHandler)
	r.Handle("/registrationForm", registrationFormHandler)
	r.Handle("/registration", registrationHandler)
	r.Handle("/authorisation", authorisationHandler)

	fmt.Println("Server is listening...")

	// Тут используется прослойка - абстрактный код, который срабатывает до
	// выполнения основного действия программы. LoggingHandler - функция, которая
	// логирует получаемые и отправляемые запросы в консоль
	http.ListenAndServe(":8181", handlers.LoggingHandler(os.Stdout, r))
}
