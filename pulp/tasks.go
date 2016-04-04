//
// Copyright 2016, Marc Sutter
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pulp

import (
	"fmt"
)

type TasksService struct {
	client *Client
}

// included in task
type Content struct {
	State        string   `json:"state"`
	ItemsTotal   int      `json:"items_total"`
	ItemsLeft    int      `json:"items_left"`
	SizeTotal    int      `json:"size_total"`
	SizeLeft     int      `json:"size_left"`
	ErrorDetails []string `json:"error_details"`
}

// included in task
type Metadata struct {
	State string
	Error string
}

// included in task
type Task struct {
	Id             string `json:"task_id"`
	StartTime      string `json:"start_time"`
	FinishTime     string `json:"finish_time"`
	State          string `json:"state"`
	Error          *Error `json:"error"`
	ProgressReport struct {

		// yum importer
		YumImporter struct {
			Content  *Content
			Metadata *Metadata
		} `json:"yum_importer"`

		// docker importer
		DockerImporter struct {
			Content  *Content
			Metadata *Metadata
		} `json:"docker_importer"`
	} `json:"progress_report"`

	Result struct {
		Details struct {
			Content *Content `json:"content"`
		} `json:"details"`
	} `json:"result"`
}

func (t Task) String() string {
	return Stringify(t)
}

func (s *TasksService) ListTasks() ([]*Task, *Response, error) {
	req, err := s.client.NewRequest("GET", "tasks/", nil)
	if err != nil {
		return nil, nil, err
	}

	var t []*Task
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

func (s *TasksService) GetTask(task string) (*Task, *Response, error) {
	u := fmt.Sprintf("tasks/%s/", task)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Task)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}
