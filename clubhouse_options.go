package main

import (
	"fmt"
	"log"

	ch "github.com/jnormington/clubhouse-go"
)

//ClubhouseOptions stores the options selected by the user
type ClubhouseOptions struct {
	Project                  *ch.Project
	State                    *ch.State
	ClubhouseEntry           *ch.Clubhouse
	StoryType                string
	AddCommentWithTrelloLink bool
	ImportMember             *ch.Member
}

// ListMember makes the call to Clubhouse package for the list
// of members. And fails hard if an err occurs.
func (co *ClubhouseOptions) ListMembers() *[]ch.Member {
	u, err := co.ClubhouseEntry.ListMembers()

	if err != nil {
		log.Fatal(err)
	}

	return &u
}

// SetupClubhouseOptions calls all the functions which consist of questions
// for building ClubhouseOptions and returns a pointer to ClubhouseOptions instance
func SetupClubhouseOptions() *ClubhouseOptions {
	var co ClubhouseOptions

	co.ClubhouseEntry = ch.New(clubHouseToken)

	co.getProjectsAndPromptUser()
	co.getWorkflowStatesAndPromptUser()
	co.getMembersAndPromptUser()
	co.promptUserForStoryType()
	co.promptUserIfAddCommentWithTrelloLink()

	return &co
}

func (co *ClubhouseOptions) promptUserIfAddCommentWithTrelloLink() {
	fmt.Println("Would you like a comment added with the original trello ticket link?")
	for i, b := range yesNoOpts {
		fmt.Printf("[%d] %s\n", i, b)
	}

	i := promptUserSelectResource()
	if i >= len(yesNoOpts) {
		log.Fatal(errOutOfRange)
	}

	if i == 0 {
		co.AddCommentWithTrelloLink = true
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

func (co *ClubhouseOptions) getMembersAndPromptUser() {
	members, err := co.ClubhouseEntry.ListMembers()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Please select a backup user account if a user is not mapped correctly")
	for i, u := range members {
		fmt.Printf("[%d] %s\n", i, u.Profile.Name)
	}

	i := promptUserSelectResource()
	if i >= len(members) {
		log.Fatal(errOutOfRange)
	}

	co.ImportMember = &members[i]
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
