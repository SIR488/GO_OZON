package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"MyMod/domain"
)

type Data struct {
	NextID domain.ID     `json:"next_id"`
	Tasks  []domain.Task `json:"tasks"`
}

type RepostoryFile struct {
	filename string
}

func (r *RepostoryFile) Init(filename string) error {
	r.filename = filename
	fileData, err := os.ReadFile(r.filename)
	if err == nil {
		var data Data
		err = json.Unmarshal(fileData, &data)
		if err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		return nil
	}
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(r.filename)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}

		jsonData, err := json.Marshal(Data{NextID: 0})
		_, err = file.Write(jsonData)
		if err != nil {
			return fmt.Errorf("write file: %w", err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("close file: %w", err)
		}
		return nil
	}
	return fmt.Errorf("open file: %w", err)
}
func (r *RepostoryFile) Save(task domain.Task) error {
	fileData, err := os.ReadFile(r.filename)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var data Data
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	task.ID = data.NextID
	data.NextID++
	data.Tasks = append(data.Tasks, task)

	out, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := os.WriteFile(r.filename, out, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
func (r *RepostoryFile) ReadAll() ([]domain.Task, error) {
	fileData, err := os.ReadFile(r.filename)
	if err != nil {
		return []domain.Task{}, fmt.Errorf("read file: %w", err)
	}

	var data Data
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return []domain.Task{}, fmt.Errorf("unmarshal: %w", err)
	}

	return data.Tasks, nil
}
func (r *RepostoryFile) ReadByID(id domain.ID) (domain.Task, error) {
	fileData, err := os.ReadFile(r.filename)
	if err != nil {
		return domain.Task{}, fmt.Errorf("read file: %w", err)
	}

	var data Data
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return domain.Task{}, fmt.Errorf("unmarshal: %w", err)
	}

	for _, task := range data.Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return domain.Task{}, domain.ErrNotFound
}
func (r *RepostoryFile) Delete(id domain.ID) error {
	fileData, err := os.ReadFile(r.filename)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var data Data
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	for i, task := range data.Tasks {
		if task.ID == id {
			data.Tasks = append(data.Tasks[:i], data.Tasks[i+1:]...)

			out, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("marshal: %w", err)
			}

			if err := os.WriteFile(r.filename, out, 0o644); err != nil {
				return fmt.Errorf("write file: %w", err)
			}
			return nil
		}
	}
	return domain.ErrNotFound
}

func (r *RepostoryFile) ChangeTask(id domain.ID, newStatus *domain.Status, newTitle *string, newDeadline *time.Time) error {
	fileData, err := os.ReadFile(r.filename)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var data Data
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	for i, task := range data.Tasks {
		if task.ID == id {
			if newStatus != nil {
				data.Tasks[i].Status = *newStatus
			}
			if newTitle != nil {
				data.Tasks[i].Title = *newTitle
			}
			if newDeadline != nil {
				data.Tasks[i].Deadline = *newDeadline
			}

			out, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("marshal: %w", err)
			}

			if err := os.WriteFile(r.filename, out, 0o644); err != nil {
				return fmt.Errorf("write file: %w", err)
			}
			return nil
		}
	}
	return domain.ErrNotFound
}
