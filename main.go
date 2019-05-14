package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

var userTodoItems = map[string]Todo{}

type Todo struct {
	Items []TodoItem `json:"items"`
}

type TodoItem struct {
	Title  string `json:"title"`
	IsDone bool   `json:"isDone"`
}

func main() {

	e := echo.New()
	e.HideBanner = true

	e.GET("/todos", func(c echo.Context) error {
		cookie := getOrGenerateCookieIfNeed(c)
		c.SetCookie(cookie)

		log.Printf("GET todo items with userID %v", cookie.Value)

		todoItems, exist := userTodoItems[cookie.Value]
		if !exist {
			todoItems = Todo{Items: []TodoItem{}}
			userTodoItems[cookie.Value] = todoItems
		}
		return c.JSON(http.StatusOK, todoItems)
	})

	e.PUT("/todos", func(c echo.Context) error {
		cookie := getOrGenerateCookieIfNeed(c)
		c.SetCookie(cookie)

		log.Printf("PUT todo items with userID %v", cookie.Value)

		var todo Todo
		if err := c.Bind(&todo); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"code":    "PUT-0001",
				"message": err.Error(),
			})
		}

		userTodoItems[cookie.Value] = todo
		return c.NoContent(http.StatusOK)
	})

	log.Fatal(e.Start(":8080"))
}

func getOrGenerateCookieIfNeed(c echo.Context) *http.Cookie {
	cookie, err := c.Cookie("user-id")
	if err != nil || cookie == nil {
		return generateNewCookie()
	}
	return cookie
}

func generateNewCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "user-id"
	cookie.Value = generateNewUserID()
	return cookie
}

func generateNewUserID() string {
	for {
		userID := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000000))
		if _, exist := userTodoItems[userID]; !exist {
			return userID
		}
	}
}
