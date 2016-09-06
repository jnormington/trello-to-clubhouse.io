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

	stdinReader   = bufio.NewReader(os.Stdin)
	errOutOfRange = "Number input is out of range. Try again"
)

func main() {
	to := setupTrelloOptionsFromUser()

	c := to.getCards()

	cards := ProcessCardsForExporting(&c)
	ExportCardsToFile(cards)

	co := setupClubhouseOptions()
	confirmAllOptionsBeforeImport(to, co)

	ImportCardsIntoClubhouse(cards, co)
	fmt.Println("*** Looks like we finished go and have fun & joy with Clubhouse ***")
}

func confirmAllOptionsBeforeImport(to *TrelloOptions, co *ClubhouseOptions) {
	opts := []string{"Yes", "No"}

	fmt.Println("****** WARNING ******")
	fmt.Println("Please review carefully before you continue")
	fmt.Printf("\nExport cards from Trello\n\tBoard: %s\n\tList: %s\n\n\n", to.Board.Name, to.List.Name)
	fmt.Printf("Import cards into clubhouse\n\tProject: %s\n\tWorkflow State: %s\n\tStory Type: %s\n\tAdd Comment with Trello Link: %t\n\n",
		co.Project.Name, co.State.Name, co.StoryType, co.AddCommentWithTrelloLink)

	fmt.Println("Is the above correct select the number representing your answer ?")

	for i, o := range opts {
		fmt.Printf("[%d] %s\n", i, o)
	}

	i := promptUserSelectResource()

	if i != 0 {
		log.Fatal("Stopping user aborted at confirmation step")
	}
}
