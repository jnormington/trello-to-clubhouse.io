package main

import (
	"fmt"
	"log"

	ch "github.com/MCBrandenburg/clubhouse-go"
)

type ClubhouseOptions struct {
	Project                 *ch.Project
	State                   *ch.State
	ClubhouseEntry          *ch.Clubhouse
	StoryType               string
	AddCommentToTrelloStory bool
}

func setupClubhouseOptions() *ClubhouseOptions {
	var co ClubhouseOptions

	co.ClubhouseEntry = ch.New(clubHouseToken)

	co.getProjectsAndPromptUser()
	co.getWorkflowStatesAndPromptUser()
	co.promptUserForStoryType()
	co.promptUserIfAddCommentWithTrelloLink()

	return &co
}

func (co *ClubhouseOptions) promptUserIfAddCommentWithTrelloLink() {
	opts := []string{"Yes", "No"}

	fmt.Println("Would you like a comment added with the original trello ticket link?")
	for i, b := range opts {
		fmt.Printf("[%d] %s\n", i, b)
	}

	i := promptUserSelectResource()
	if i >= len(opts) {
		log.Fatal(errOutOfRange)
	}

	if i == 0 {
		co.AddCommentToTrelloStory = true
	}
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
		log.Fatal(errOutOfRange)
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
