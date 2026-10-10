package ui

import (
	"MyMod/domain"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	Save(title string, deadline time.Time) error
	ReadAll() ([]domain.Task, error)
	ReadByID(id domain.ID) (domain.Task, error)
	Delete(id domain.ID) error
	ChangeStatus(id domain.ID, newStatus domain.Status) error
	ChangeTask(id domain.ID, newTitle *string, newDeadline *time.Time) error
}
type UI struct {
	Service Service
}

func (ui *UI) Start() {
	fmt.Println("start working")
	fmt.Println("possible commands: help, exit, create, read all, read one, change status, change task, delete")

	scanner := bufio.NewScanner(os.Stdin)
	var title string
	var status domain.Status
	var deadlineStr string
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		fmt.Print("type command")
		message := strings.TrimSpace(scanner.Text())
		switch message {
		case "help":
			fmt.Println("possible commands: help, exit, create, read all, read one, change status, change task, delete")
		case "exit":
			return
		case "create":
			fmt.Println("write title")
			if !scanner.Scan() {
				return
			}
			title = scanner.Text()
			fmt.Println("write deadline YYYY-MM-DD")
			if !scanner.Scan() {
				return
			}
			deadlineStr = scanner.Text()
			deadline, err := time.Parse(time.DateOnly, deadlineStr)
			if err != nil {
				fmt.Println("invalid deadline:", err)
				continue
			}
			ui.CreateTask(title, deadline)
		case "read all":
			ui.ReadAll()
		case "read one":
			fmt.Println("write ID")
			if !scanner.Scan() {
				return
			}
			id, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("invalid id:", err)
				continue
			}
			ui.ReadByID(domain.ID(id))
		case "change status":
			fmt.Println("write ID")
			if !scanner.Scan() {
				return
			}
			id, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("invalid id:", err)
				continue
			}
			fmt.Println("write number of new status")
			fmt.Printf("0: %s\n1: %s\n2: %s\n3: %s\n",
				domain.StatusPlanned, domain.StatusInProgress, domain.StatusCanceled, domain.StatusDone)
			if !scanner.Scan() {
				return
			}
			num, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("invalid number:", err)
				continue
			}
			switch num {
			case 0:
				status = domain.StatusPlanned
			case 1:
				status = domain.StatusInProgress
			case 2:
				status = domain.StatusCanceled
			case 3:
				status = domain.StatusDone
			default:
				fmt.Println("unknown status")
				continue
			}
			ui.ChangeStatus(domain.ID(id), status)
		case "change task":
			fmt.Println("write ID")
			if !scanner.Scan() {
				return
			}
			id, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("invalid id:", err)
				continue
			}
			fmt.Println("write new title (or press enter)")
			if !scanner.Scan() {
				return
			}
			title = scanner.Text()
			fmt.Println("write new deadline YYYY-MM-DD (or press enter)")
			if !scanner.Scan() {
				return
			}
			deadlineStr = scanner.Text()
			deadline, err := time.Parse(time.DateOnly, deadlineStr)
			if err != nil {
				ui.ChangeTask(domain.ID(id), title, nil)
			} else {
				ui.ChangeTask(domain.ID(id), title, &deadline)
			}
		case "delete":
			fmt.Println("write ID")
			if !scanner.Scan() {
				return
			}
			id, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("invalid id:", err)
				continue
			}
			ui.Delete(domain.ID(id))
		default:
			fmt.Println("unknown command, try 'help'")
			continue
		}
		if scanner.Err() != nil {
			fmt.Println("input error")
			time.Sleep(3 * time.Second)
			return
		}
	}
}

func (ui *UI) CreateTask(title string, deadline time.Time) {
	err := ui.Service.Save(title, deadline)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Println("succesfull create")
}

func (ui *UI) ReadAll() {
	tasks, err := ui.Service.ReadAll()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	fmt.Println("ID  status  title  deadline")
	for _, task := range tasks {
		fmt.Printf("%d  %s  '%s'  %s\n", task.ID, task.Status, task.Title, task.Deadline.Format(time.DateOnly))
	}
}
func (ui *UI) ReadByID(id domain.ID) {
	task, err := ui.Service.ReadByID(id)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Println("ID  status  title  deadline")
	fmt.Printf("%d  %s  '%s'  %s\n", task.ID, task.Status, task.Title, task.Deadline.Format(time.DateOnly))

}
func (ui *UI) Delete(id domain.ID) {
	err := ui.Service.Delete(id)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Println("succesfull remove")
}

func (ui *UI) ChangeStatus(id domain.ID, newStatus domain.Status) {
	err := ui.Service.ChangeStatus(id, newStatus)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Println("succesfull change")
}
func (ui *UI) ChangeTask(id domain.ID, newTitle string, newDeadline *time.Time) {
	var title *string
	if newTitle != "" {
		title = &newTitle
	}
	err := ui.Service.ChangeTask(id, title, newDeadline)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Println("succesfull change")
}
