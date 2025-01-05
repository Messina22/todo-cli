package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Task represents a task in the ToDo list
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// TaskList represents a list of tasks
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

const fileName = ".tasks.json"

func LoadTasks() TaskList {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return TaskList{}
	}
	tasks := TaskList{}

	if err := json.Unmarshal(file, &tasks); err != nil {
		fmt.Println("Error loading tasks: ", err)
		return TaskList{}
	}

	return tasks
}

func SaveTasks(tasks TaskList) {
	file, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Error saving tasks: ", err)
		return
	}

	if err := os.WriteFile(fileName, file, 0666); err != nil {
		fmt.Println("Error saving tasks: ", err)
	}
}

// CLI ToDo List
func main() {
	taskList := LoadTasks()

	// Define flags
	add := flag.String("add", "", "Add a new task")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Mark a task complete by its ID")
	remove := flag.Int("remove", 0, "Remove a task by its ID")
	flag.Parse()

	if *add != "" {
		// Add a new task
		newTask := Task{
			ID:          len(taskList.Tasks) + 1,
			Description: *add,
			Completed:   false,
		}
		taskList.Tasks = append(taskList.Tasks, newTask)
		SaveTasks(taskList)
		fmt.Println("Adding task: ", newTask.Description)
	} else if *list {
		// List all tasks
		for _, task := range taskList.Tasks {
			status := "Pending"
			if task.Completed {
				status = "Completed"
			}
			fmt.Printf("%d: %s (%s)\n", task.ID, task.Description, status)
		}
	} else if *complete > 0 {
		// Mark a task as completed
		for i, task := range taskList.Tasks {
			if task.ID == *complete {
				taskList.Tasks[i].Completed = true
				SaveTasks(taskList)
				fmt.Println("Task completed: ", task.Description)
				return
			}
		}
		fmt.Println("Task not found")
	} else if *remove > 0 {

		// Remove a task
		newTasks := []Task{}

		for _, task := range taskList.Tasks {
			i := 1
			if task.ID != *remove {
				task.ID = i
				newTasks = append(newTasks, task)
				i += 1
			}
		}

		if len(newTasks) == len(taskList.Tasks) {
			fmt.Println("Task not found")
		} else {
			taskList.Tasks = newTasks
			SaveTasks(taskList)
			fmt.Println("Task removed")
		}
	} else {
		fmt.Println("Usage: ")
		flag.PrintDefaults()
	}
}
