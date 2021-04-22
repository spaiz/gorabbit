//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"gorabbit/infra/cli/pkg"
	"log"
)

// -----------------------------------------------------------//

// Rescale will regenerate mappings, shut down and up the instances again
func Rescale() error {
	mg.Deps(Remap, ShutDown, Up)
	return nil
}

// Remap generates new nomad jobs with new mapping generates as result of manifest.json
func Remap() error {
	manifest, err := pkg.LoadManifest("manifest.json")
	if err != nil {
		return err
	}

	state, err := manifest.CreateState()
	if err != nil {
		return err
	}

	jobs, err := manifest.CreateJobs(state)
	if err != nil {
		return err
	}

	return manifest.DumpJobs(jobs, "../nomad/jobs/")
}

// -----------------------------------------------------------//
// SystemUp will startup RabbitMQ
func SystemUp() error {
	mg.Deps(RabbitMQUp)
	return nil
}

func SystemShutDown() error {
	mg.Deps(ShutDown, RabbitMQShutDown)
	return nil
}

func Up() error {
	mg.Deps(Remap, ConsumerUp)
	return nil
}

func ShutDown() error {
	mg.Deps(ConsumerShutDown)
	return nil
}

func ConsumerUp() error {
	manifest, err := pkg.LoadManifest("manifest.json")
	if err != nil {
		return err
	}

	for _, worker := range manifest.Workers {
		job := fmt.Sprintf("../nomad/jobs/%s.nomad", worker.Identity)
		log.Printf("Ranning job: %s", job)
		err := sh.Run("nomad", "job", "run", job)
		if err != nil {
			log.Printf("Failed to start job: %s (%s)", worker.Identity, err)
		}
	}

	return nil
}

func ConsumerShutDown() error {
	manifest, err := pkg.LoadManifest("manifest.json")
	if err != nil {
		return err
	}

	for _, worker := range manifest.Workers {
		err := sh.Run("nomad", "job", "stop", "-purge", worker.Identity)
		if err != nil {
			log.Printf("Failed to stop job: %s (%s)", worker.Identity, err)
		}
	}

	return nil
}

func RabbitMQUp() error {
	return sh.Run("nomad", "job", "run", "../nomad/jobs/rabbitmq.nomad")
}

func RabbitMQShutDown() error {
	return sh.Run("nomad", "job", "stop", "-purge", "rabbitmq")
}

// -----------------------------------------------------------//
// InstallConsumer builds and installs the binary of producer
func InstallConsumer() error {
	return sh.Run("go", "install", "../../cmd/producer")
}

// InstallConsumer builds and installs the binary of consumer
func InstallProducer() error {
	return sh.Run("go", "install", "../../cmd/consumer")
}

// Install builds and installs the binary of consumer and producer
func Install() error {
	mg.Deps(InstallProducer, InstallConsumer)
	return nil
}
