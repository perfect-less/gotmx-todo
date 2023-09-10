package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"
)

const PORT string = "8080"
const IDLENGTH = 10

var todos []Todo
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//go:embed web/static/css/style.css web/index.html web/htmx.min.js
var content embed.FS

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

func handleSave(w http.ResponseWriter, r *http.Request) {
	jBytes, err := json.MarshalIndent(todos, "", " ")
	if err != nil {
		fmt.Printf("err on json.MarshallIndent: %v\n", err)
	}

	err = os.WriteFile("todos.save", jBytes, 0644)
	if err != nil {
		fmt.Printf("err on os.WriteFile: %v\n", err)
	}

	time.Sleep(1 * time.Second)
}

func loadSave() bool {
	_, err := os.Stat("todos.save")
	if err == nil {
		jsonFile, err := os.Open("todos.save")
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return false
		}

		jsonBytes, err := io.ReadAll(jsonFile)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return false
		}

		err = json.Unmarshal(jsonBytes, &todos)
		if err != nil {
			return false
		}
	} else {
		return false
	}

	return true
}

func main() {
	fmt.Println("Running gotmx")

	if !loadSave() {
		addNewTodo("Install gotmx-todo", "You've downloaded and installed gotmx-todo on your device, click the check button on the right to mark task as done")
		addNewTodo("Write essays", "Don't forget to write your essay assignment for CS101 course.")
	}

	http.HandleFunc("/", handleLanding)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/check-done", handleCheckDone)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/save", handleSave)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(content))))

	fmt.Printf("Listening to: http://localhost:%s\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}
