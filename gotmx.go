package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

type Todo struct {
	Id          string
	Todo        string
	Description string
	Done        bool
}

const PORT string = "8080"
const IDLENGTH = 10

var todos []Todo

//go:embed web/static/css/style.css web/index.html web/htmx.min.js
var content embed.FS

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func addNewTodo(todo string, description string) {
	newId := ""
	for {
		newId = randSeq(IDLENGTH)

		isUnique := true
		for i := range todos {
			if todos[i].Id == newId {
				isUnique = false
				break
			}
		}
		if isUnique {
			break
		}
	}

	todos = append(todos, Todo{Id: newId, Todo: todo, Description: description})
}

func removeByIndex(slice []Todo, s int) []Todo {
	return append(slice[:s], slice[s+1:]...)
}

func executeTodoListTemplate(w http.ResponseWriter) error {
	temp := template.Must(template.ParseFS(content, "web/index.html"))
	todosData := map[string][]Todo{
		"todos": todos,
	}

	err := temp.ExecuteTemplate(w, "todo-list", todosData)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		fmt.Fprintf(w, "<h1>Error code %d encountered<h1>", http.StatusInternalServerError)
	}
	return err
}

func handleLanding(w http.ResponseWriter, r *http.Request) {
	var temp = template.Must(template.ParseFS(content, "web/index.html"))
	todosData := map[string][]Todo{
		"todos": todos,
	}

	err := temp.Execute(w, todosData)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		fmt.Fprintf(w, "<h1>Error code %d encountered<h1>", http.StatusInternalServerError)
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	todo := r.PostFormValue("todo")
	description := r.PostFormValue("description")

	addNewTodo(todo, description)
	executeTodoListTemplate(w)
}

func handleCheckDone(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Failed reading rquest body: %s\n", err)
	}

	bodyStr, err := url.QueryUnescape(string(resBody))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	bodySplited := strings.Split(bodyStr, "=")
	todoId := bodySplited[len(bodySplited)-1]
	fmt.Printf("bodySplited[len(bodySplited)-1]: %v\n", bodySplited[len(bodySplited)-1])

	for i := range todos {
		if todoId == todos[i].Id {
			todos[i].Done = true
			break
		}
	}

	executeTodoListTemplate(w)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Failed reading rquest body: %s\n", err)
	}

	bodyStr, err := url.QueryUnescape(string(resBody))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	bodySplited := strings.Split(bodyStr, "=")
	todoId := bodySplited[len(bodySplited)-1]
	fmt.Printf("bodySplited[len(bodySplited)-1]: %v\n", bodySplited[len(bodySplited)-1])

	i := 0
	for i = range todos {
		if todoId == todos[i].Id {
			break
		}
	}

	todos = removeByIndex(todos, i)
	executeTodoListTemplate(w)
}

func main() {
	fmt.Println("Running gotmx")

	addNewTodo("Install gotmx-todo", "You've downloaded and installed gotmx-todo on your device, click the check button on the right to mark task as done")
	addNewTodo("Write essays", "Don't forget to write your essay assignment for CS101 course.")

	http.HandleFunc("/", handleLanding)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/check-done", handleCheckDone)
	http.HandleFunc("/delete", handleDelete)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(content))))

	fmt.Printf("Listening to: http://localhost:%s\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}
