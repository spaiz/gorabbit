package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

type Manifest struct {
	Workers []*Worker `json:"workers"`
}

type State struct {
	mappings []*QueuesMapping
}

func (s *State) addMapping(mapping *QueuesMapping) {
	s.mappings = append(s.mappings, mapping)
}

type WorkerInstance struct {
	Identity                 string
	InstanceId               int
	queues                   []string
	RabbitMQConnectionString string
}

func (r *WorkerInstance) addQueue(name string) {
	r.queues = append(r.queues, name)
}

func (m *Manifest) CreateState() (*State, error) {
	state := &State{
		mappings: []*QueuesMapping{},
	}

	for _, worker := range m.Workers {
		mapping := m.createMapping(worker)
		state.addMapping(mapping)
	}

	return state, nil
}

type QueuesMapping struct {
	instances []*WorkerInstance
	topic     string
	identity  string
}

func (m *Manifest) createMapping(worker *Worker) *QueuesMapping {
	instances := make([]*WorkerInstance, 0, worker.InstancesNumber)
	mapping := &QueuesMapping{
		identity: worker.Identity,
		topic:    worker.Topic,
	}

	for i := 0; i < worker.InstancesNumber; i++ {
		instance := &WorkerInstance{
			Identity:                 worker.Identity,
			InstanceId:               i + 1,
			queues:                   []string{},
			RabbitMQConnectionString: worker.RabbitMQConnectionString,
		}

		instances = append(instances, instance)
	}

	for i := 1; i <= worker.QueuesSet.PartitionsNumber; i++ {
		queueName := fmt.Sprintf("%s_%04d", worker.Identity, i)
		workerId := i % worker.InstancesNumber
		instances[workerId].addQueue(queueName)
	}

	mapping.instances = instances
	return mapping
}

func (m *Manifest) Validate() error {
	for _, worker := range m.Workers {
		err := worker.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manifest) CreateJobs(state *State) ([]*Job, error) {
	jobs := make([]*Job, 0, 0)

	for _, mapping := range state.mappings {
		args := []string{""}
		argsFormatted, _ := json.Marshal(args)

		job := &Job{
			Name: mapping.identity,
			Group: &Group{
				Tasks: make([]*Task, 0, 0),
				Name:  fmt.Sprintf("%s_%s", mapping.topic, mapping.identity),
			},
		}

		for i, instance := range mapping.instances {
			task := &Task{
				Name:                     fmt.Sprintf("%s_%04d", instance.Identity, i),
				Queues:                   strings.Join(instance.queues, " "),
				Identity:                 mapping.identity,
				Topic:                    mapping.topic,
				QueueBindingKey:          "1",
				Cmd:                      "consumer",
				Args:                     string(argsFormatted),
				RabbitMQConnectionString: instance.RabbitMQConnectionString,
			}

			job.Group.AddTask(task)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

type Job struct {
	Name  string
	Group *Group
}

func (m *Manifest) DumpJobs(jobs []*Job, dir string) error {
	jobTmplFilepath := "nomad.tmpl"
	t, err := template.ParseFiles(jobTmplFilepath)
	if err != nil {
		log.Printf("Failed to parse nomad.tmpl (%s)", err)
		return err
	}

	for _, job := range jobs {
		buf := bytes.Buffer{}
		err = t.Execute(&buf, job)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/%s.nomad", dir, job.Name)
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}

		_, err = file.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}

type Worker struct {
	Identity                 string    `json:"identity"`
	Topic                    string    `json:"topic"`
	QueuesSet                QueuesSet `json:"queuesSet"`
	InstancesNumber          int
	RabbitMQConnectionString string `json:"rabbitMQConnectionString"`
}

func (w *Worker) Validate() error {
	if w.InstancesNumber > w.QueuesSet.PartitionsNumber {
		return fmt.Errorf("you cannot have more consumers then queues: [Identity] %s, [InstancesNumber] %d, [PartitionsNumber] %d", w.Identity, w.InstancesNumber, w.QueuesSet.PartitionsNumber)
	}

	return nil
}

type QueuesSet struct {
	PartitionsNumber int `json:"partitionsNumber"`
}

func LoadManifest(filePath string) (*Manifest, error) {
	manifest := &Manifest{}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &Manifest{}, err
	}

	err = json.Unmarshal(content, manifest)
	if err != nil {
		log.Printf("Failed to unmarshal manifest file: %s", err)
	}

	err = manifest.Validate()
	return manifest, err
}

type Group struct {
	Tasks []*Task
	Name  string
}

func (g *Group) AddTask(task *Task) {
	g.Tasks = append(g.Tasks, task)
}

type Task struct {
	Queues                   string
	QueueBindingKey          string
	Cmd                      string
	Args                     string
	Name                     string
	RabbitMQConnectionString string
	Topic                    string
	Identity                 string
}
