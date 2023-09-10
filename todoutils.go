package main

import "math/rand"

type Todo struct {
	Id          string
	Todo        string
	Description string
	Done        bool
}

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
