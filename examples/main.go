// simple example of getting a repo metadatas and running a sync on it
package main

import (
	"fmt"
	"github.com/msutter/go-pulp/pulp"
	"log"
	"time"
)

func main() {
	apiUser := "admin"
	apiPasswd := "admin"
	apiEndpoint := "pulp-lab-11.test"

	DisableSsl := false
	SkipSslVerify := true

	// create the client
	client, err := pulp.NewClient(apiEndpoint, apiUser, apiPasswd, DisableSsl, SkipSslVerify, nil)

	// repository options
	ro := &pulp.GetRepositoryOptions{
		Details: true,
	}

	repo := "sccloud-mgmt-infra-el6-lab"

	// get the repo
	r, _, rerr := client.Repositories.GetRepository(repo, ro)
	fmt.Printf("%v\n", r)

	if rerr != nil {
		fmt.Println(rerr.Error())
		log.Fatal(rerr)
	}

	// sync it
	syncCallReport, _, err := client.Repositories.SyncRepository(repo)
	syncTaskId := syncCallReport.SpawnedTasks[0].TaskId
	fmt.Printf("TaskId: %v\n", syncTaskId)
	if err != nil {
		log.Fatal(err)
	}

	state := "init"
	for (state != "finished") && (state != "error") {
		task, _, terr := client.Tasks.GetTask(syncTaskId)

		fmt.Printf("----- progress --------\n")
		fmt.Printf("state: %v\n", task.State)
		fmt.Printf("progressReport: %v\n", task.ProgressReport)

		var importer *pulp.Importer
		if task.Importer() == "yum" {
			importer = task.ProgressReport.YumImporter
		}
		if task.Importer() == "docker" {
			importer = task.ProgressReport.DockerImporter
		}

		fmt.Printf("importer: %v\n", task.Importer())
		fmt.Printf("item Total: %v\n", importer.Content.ItemsTotal)
		fmt.Printf("item Left: %v\n", importer.Content.ItemsLeft)
		state = task.State
		time.Sleep(500 * time.Millisecond)
		if terr != nil {
			log.Fatal(terr)
		}
	}
}
