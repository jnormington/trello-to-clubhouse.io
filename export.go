package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path"

	trello "github.com/Luzifer/go-trello"
)

type Card struct {
	Card     trello.Card       `json:"Card"`
	Creator  string            `json:"card_creator"`
	Comments map[string]string `json:"comments"`
}

func GetCommentsForCards(crds *[]trello.Card) *[]Card {
	var cards []Card

	for _, card := range *crds {
		var c Card

		c.Card = card

		actions, err := card.Actions()
		if err != nil {
			fmt.Println("Error: Querying the actions for:", card.Name, "failed", "continuing..", err)
		}

		for _, a := range actions {
			if a.Type == "commentCard" && a.Data.Text != "" {
				c.Comments = map[string]string{
					"Creator": a.MemberCreator.FullName,
					"Text":    a.Data.Text,
				}
			} else if a.Type == "createCard" {
				c.Creator = a.MemberCreator.FullName
			}
		}

		cards = append(cards, c)
	}

	return &cards
}

func ExportCardsToFile(cards *[]Card) {
	b, err := json.Marshal(cards)
	if err != nil {
		log.Fatal(err)
	}

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	filePath := path.Join(u.HomeDir, "cards.json")
	ioutil.WriteFile(filePath, b, 0644)
}
