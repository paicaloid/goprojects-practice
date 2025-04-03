package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
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
