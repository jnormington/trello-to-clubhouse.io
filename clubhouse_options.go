package main

import (
	"fmt"
	"log"

	ch "github.com/MCBrandenburg/clubhouse-go"
)

type ClubhouseOptions struct {
	Project        *ch.Project
	State          *ch.State
	ClubhouseEntry *ch.Clubhouse
	StoryType      string
}

func setupClubhouseOptions() *ClubhouseOptions {
	var co ClubhouseOptions

	co.ClubhouseEntry = ch.New(clubHouseToken)

	co.getProjectsAndPromptUser()
	co.getWorkflowStatesAndPromptUser()
	co.promptUserForStoryType()

	return &co
}

func (co *ClubhouseOptions) getProjectsAndPromptUser() {
	projects, err := co.ClubhouseEntry.ListProjects()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Please select a project by it number to import the cards into")
	for i, p := range projects {
		fmt.Printf("[%d] %s\n", i, p.Name)
	}

	i := promptUserSelectResource()
	if i >= len(projects) {
		log.Fatal("Number input is out of range ?")
	}

	co.Project = &projects[i]
}

func (co *ClubhouseOptions) getWorkflowStatesAndPromptUser() {
	workflows, err := co.ClubhouseEntry.ListWorkflow()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Please select a workflow state to import the cards into")
	var statesLen int

	for _, w := range workflows {
		statesLen += len(w.States)

		for i, s := range w.States {
			fmt.Printf("[%d] %s\n", i, s.Name)
		}
	}

	i := promptUserSelectResource()
	if i >= statesLen {
		log.Fatal(errOutOfRange)
	}

	co.State = &workflows[0].States[i]
}

func (co *ClubhouseOptions) promptUserForStoryType() {
	types := []string{"feature", "chore", "bug"}

	fmt.Println("Please select the story type all cards should be imported as")
	for i, t := range types {
		fmt.Printf("[%d] %s\n", i, t)
	}

	i := promptUserSelectResource()
	if i >= len(types) {
		log.Fatal(errOutOfRange)
	}

	co.StoryType = types[i]
}
