package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	clubHouseToken = os.Getenv("CLUBHOUSE_TOKEN")
	trelloToken    = os.Getenv("TRELLO_TOKEN")
	trelloKey      = os.Getenv("TRELLO_KEY")
	dropboxToken   = os.Getenv("DROPBOX_TOKEN")

	stdinReader   = bufio.NewReader(os.Stdin)
	errOutOfRange = "Number input is out of range. Try again"
	yesNoOpts     = []string{"Yes", "No"}
)

func main() {
	to := SetupTrelloOptionsFromUser()

	c := to.getCards()

	cards := ProcessCardsForExporting(&c, to)

	co := SetupClubhouseOptions()

	mapTrelloToClubhouseUsers(to, co)
	confirmAllOptionsBeforeImport(to, co)

	ImportCardsIntoClubhouse(cards, co)
	fmt.Println("*** Looks like we finished go and have fun & joy with Clubhouse ***")
}

func mapTrelloToClubhouseUsers(to *TrelloOptions, co *ClubhouseOptions) {
	tMembers, err := to.Board.Members()
	co.UserMapping = make(map[string]string)
	missing := false

	if err != nil {
		log.Fatal(err)
	}

	cUsers := co.ListUsers()

	for _, m := range tMembers {
		match := false

		for _, u := range cUsers {
			// Email address interestingly sits within the permission
			// Validate that we have at least one permission entry
			if len(u.Permissions) == 0 {
				continue
			}

			if m.Email == u.Permissions[0].EmailAddress {
				match = true
				co.UserMapping[m.Id] = u.ID
			}
		}

		if !match {
			fmt.Printf("User %s not found in clubhouse\n", m.Email)
			missing = true
		}
	}

	// Now check if any are missing and prompt to default them to one user
	if missing {
		fmt.Printf("The above account are missing.\n")
		fmt.Printf("Please select an account to default ticket creators of these missing accounts")
		co.GetBackupUserForImport(cUsers)
	}
}

func confirmAllOptionsBeforeImport(to *TrelloOptions, co *ClubhouseOptions) {
	fmt.Println("****** WARNING ******")
	fmt.Println("Please review carefully before you continue")
	fmt.Printf("\nExport cards from Trello\n\tBoard: %s\n\tList: %s\n\n\n", to.Board.Name, to.List.Name)
	fmt.Printf("Import cards into clubhouse\n\tProject: %s\n\tWorkflow State: %s\n\tStory Type: %s\n\tAdd Comment with Trello Link: %t\n\n",
		co.Project.Name, co.State.Name, co.StoryType, co.AddCommentWithTrelloLink)

	fmt.Println("Is the above correct select the number representing your answer ?")

	for i, o := range yesNoOpts {
		fmt.Printf("[%d] %s\n", i, o)
	}

	i := promptUserSelectResource()

	if i != 0 {
		log.Fatal("Stopping user aborted at confirmation step")
	}
}
