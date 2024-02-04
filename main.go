package main

import (
  "net/http"
  "os"
  "time"

  "gorm.io/gorm"
  "gorm.io/driver/postgres"
  "gorm.io/driver/sqlite"

  "github.com/gin-gonic/gin"
)

type Todo struct {
  gorm.Model
  ID   uint   `gorm:"primaryKey"`
  Text string `gorm:"name:text"`
  Done bool   `gorm:"name:done"`
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func main() {
  db := connectDB()
  db.AutoMigrate(&Todo{})

  r := gin.Default()
  r.LoadHTMLGlob("templates/*")
  r.Static("/assets", "./assets")

  r.GET("/api/todo", getTodos(db))
  r.POST("/api/todo", putTodo(db))
  r.POST("/api/todo/:id/complete", completeTodo(db))
  r.Run(":8080")
}

func connectDB() *gorm.DB {
  var db *gorm.DB
  var err error
  var dialector gorm.Dialector

  dbType := os.Getenv("DB_ENGINE")
  
  remainingTries := 5
  for remainingTries > 0 {
    if dbType == "POSTGRES" {
      dsn := getEnv("DB_DSN", "")
      dialector = postgres.Open(dsn)
    } else {
      filename := getEnv("DB_FILE", "todo.db")
      dialector = sqlite.Open(filename)
    }

    db, err = gorm.Open(dialector, &gorm.Config{})
    if err != nil {
      time.Sleep(5 * time.Second)
      remainingTries--
    } else {
      break
    }
  }

  if err != nil {
    panic("Failed to connect to database")
  }

  return db
}

func findNonCompleteTodos(db *gorm.DB) []Todo {
  var todos []Todo
  db.Find(&todos, map[string]interface{}{"done": false})
  return todos
}

func findAllTodos(db *gorm.DB) []Todo {
  var todos []Todo
  db.Find(&todos)
  return todos
}

func isWithDone(c *gin.Context) bool {
  doneParam := c.DefaultQuery("done", "off")
  if doneParam == "on" {
    return true
  } else {
    return false
  }
}

func getTodos(db *gorm.DB) func(*gin.Context) {
  return func(c *gin.Context) {
    var todos []Todo
    if isWithDone(c) {
      todos = findAllTodos(db)
    } else {
      todos = findNonCompleteTodos(db)
    }
    c.HTML(http.StatusOK, "todos.tmpl", gin.H{
      "todos": todos,
      "withDone": isWithDone(c),
    })
  }
}

func putTodo(db *gorm.DB) func(*gin.Context) {
  return func(c *gin.Context) {
    text := c.PostForm("text")
    todo := Todo{Text: text, Done: false}
    result := db.Create(&todo)
    if result.Error != nil {
      c.AbortWithStatus(http.StatusInternalServerError)
    }
    getTodos(db)(c)
  }
}

func completeTodo(db *gorm.DB) func(*gin.Context) {
  return func(c *gin.Context) {
    id := c.Param("id")
    todo := Todo{}
    result := db.First(&todo, id)
    if result.Error != nil || result.RowsAffected == 0 {
      c.AbortWithStatus(http.StatusInternalServerError)
    }
    todo.Done = true
    db.Save(&todo)
    getTodos(db)(c)
  }
}
