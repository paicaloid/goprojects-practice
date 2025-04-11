package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
)

type TaskData struct {
	id      int
	task    string
	created string
	done    bool
}

var layout string = "2006-01-02 15:04:05 MST"

func rowToTaskData(row []string) (TaskData, error) {
	id, err := strconv.Atoi(row[0])
	if err != nil {
		return TaskData{}, err
	}

	created, err := time.Parse(layout, row[2])
	if err != nil {
		return TaskData{}, err
	}

	status, err := strconv.ParseBool(row[3])
	if err != nil {
		return TaskData{}, err
	}

	return TaskData{
		id:      id,
		task:    row[1],
		created: timediff.TimeDiff(created),
		done:    status,
	}, nil

}

func LoadFile() (*os.File, error) {
	file, err := os.OpenFile("data.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WriteFile(file *os.File, task string) {
	writer := csv.NewWriter(file)
	defer writer.Flush()

	strTime := time.Now().Format(layout)
	id, err := getLastID(file)
	if err != nil {
		log.Panic(err)
	}

	record := []string{
		strconv.Itoa(id + 1), task, strTime, "false",
	}
	writer.Write(record)
}

func ReadFile(file *os.File) ([]TaskData, error) {
	tasks := make([]TaskData, 0)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return tasks, err
	}

	for _, record := range records {
		data, err := rowToTaskData(record)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, data)
	}
	return tasks, nil
}

func ShowList(tasks []TaskData, flag bool) {
	w := new(tabwriter.Writer)
	// w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	w.Init(os.Stdout, 0, 8, 3, '\t', 0)

	if flag {
		fmt.Fprintln(w, "ID\tTask\tCreated\tDone\t")
	} else {
		fmt.Fprintln(w, "ID\tTask\tCreated\t")
	}

	for _, task := range tasks {
		if flag {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t\n", task.id, task.task, task.created, task.done)
		} else {
			if !task.done {
				fmt.Fprintf(w, "%v\t%v\t%v\t\n", task.id, task.task, task.created)
			}
		}
	}

	fmt.Fprintln(w)
	w.Flush()
}

func getLastID(file *os.File) (int, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, err
	}

	// case no data
	if len(records) == 0 {
		return 0, nil
	} else { // get id from lastrow
		data, err := rowToTaskData(records[len(records)-1])
		if err != nil {
			return 0, err
		}
		return data.id, nil
	}
}
