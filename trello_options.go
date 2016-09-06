package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	trello "github.com/jnormington/go-trello"
)

// TrelloOptions holds the links to the board and list
// from the trello api calls that the user has selected to import
type TrelloOptions struct {
	Board *trello.Board
	List  *trello.List
	User  *trello.Member
}

func setupTrelloOptionsFromUser() *TrelloOptions {
	var t TrelloOptions

	t.getCurrentUser()
	t.getBoardsAndPromptUser()
	t.getListsAndPromptUser()

	return &t
}

func (t *TrelloOptions) getCurrentUser() {
	c, err := trello.NewAuthClient(trelloKey, &trelloToken)
	if err != nil {
		log.Fatal(err)
	}

	u, err := c.Member("me")
	if err != nil {
		log.Fatal(err)
	}

	t.User = u
}

func (t *TrelloOptions) getBoardsAndPromptUser() {
	boards, err := t.User.Boards()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Please select a board by its number")
	for i, b := range boards {
		fmt.Printf("[%d] %s\n", i, b.Name)
	}

	i := promptUserSelectResource()
	if i >= len(boards) {
		log.Fatal(errOutOfRange)
	}

	t.Board = &boards[i]
}

func (t *TrelloOptions) getListsAndPromptUser() {
	lists, err := t.Board.Lists()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Please select the list to import by number")
	for i, l := range lists {
		fmt.Printf("[%d] %s\n", i, l.Name)
	}

	i := promptUserSelectResource()
	if i >= len(lists) {
		log.Fatal(errOutOfRange)
	}

	t.List = &lists[i]
}

func (t TrelloOptions) getCards() []trello.Card {
	fmt.Println("Please wait while we retrieve your cards... This might take a few minutes.")

	cards, err := t.List.Cards()
	if err != nil {
		log.Fatal(err)
	}

	return cards
}

func promptUserSelectResource() int {
	i, err := stdinReader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	i = strings.TrimRight(i, "\n")
	id, err := strconv.Atoi(i)
	if err != nil {
		log.Fatal("Hmm... did you type a number from the list ?", err)
	}

	return id
}
