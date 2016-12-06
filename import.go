package main

import (
	"fmt"
	"time"

	ch "github.com/jnormington/clubhouse-go"
)

var outputFormat = "%-40s %-17s %s\n"

// ImportCardsIntoClubhouse takes *[]Card, *ClubhouseOptions and builds a clubhouse Story
// this story from both the card and clubhouse options and creates via the api.
func ImportCardsIntoClubhouse(cards *[]Card, opts *ClubhouseOptions, um *UserMap) {
	fmt.Println("Importing trello cards into Clubhouse...")
	fmt.Printf(outputFormat+"\n", "Trello Card Link", "Import Status", "Error/Story ID")

	for _, c := range *cards {
		//We could use bulk update but lets give the user some prompt feedback
		st, err := opts.ClubhouseEntry.CreateStory(*buildClubhouseStory(&c, opts, um))
		if err != nil {
			fmt.Printf(outputFormat, c.ShortURL, "Failed", err)
			continue
		}

		fmt.Printf(outputFormat, c.ShortURL, "Success", fmt.Sprintf("Story ID: %d", st.ID))
	}
}

func buildLinkFiles(card *Card, opts *ClubhouseOptions) []int64 {
	var ids []int64

	for _, si := range card.Attachments {
		lf := ch.CreateLinkedFile{
			Name:         si.Name,
			Type:         "url",
			URL:          si.URL,
			ThumbnailURL: si.URL,
			ContentType:  si.MimeType,
			UploaderID:   opts.ImportUser.ID,
		}

		r, err := opts.ClubhouseEntry.CreateLinkedFiles(lf)
		if err != nil {
			fmt.Println("Fail to create linked file card name:", card.Name, "Dropbox link:", si.URL, "Err:", err)
		} else {
			ids = append(ids, r.ID)
		}
	}

	return ids
}

func buildClubhouseStory(card *Card, opts *ClubhouseOptions, um *UserMap) *ch.CreateStory {

	return &ch.CreateStory{
		ProjectID:       opts.Project.ID,
		WorkflowStateID: opts.State.ID,
		RequestedByID:   um.GetCreator(card.IDCreator),
		OwnerIds:        mapOwnersFromTrelloCard(card, um),
		StoryType:       opts.StoryType,

		Name:        card.Name,
		Description: card.Desc,
		Deadline:    card.DueDate,
		CreatedAt:   card.CreatedAt,

		Labels:   *buildLabels(card),
		Tasks:    *buildTasks(card),
		Comments: *buildComments(card, opts.AddCommentWithTrelloLink, um),

		LinkedFileIds: buildLinkFiles(card, opts),
	}
}

func mapOwnersFromTrelloCard(c *Card, um *UserMap) []string {
	var owners []string

	for _, o := range c.IDOwners {
		owners = append(owners, um.GetCreator(o))
	}

	return owners
}

func buildComments(card *Card, addCommentWithTrelloLink bool, um *UserMap) *[]ch.CreateComment {
	var comments []ch.CreateComment

	for _, cm := range card.Comments {
		com := ch.CreateComment{
			CreatedAt: *cm.CreatedAt,
			AuthorID:  um.GetCreator(cm.IDCreator),
			Text:      cm.Text,
		}

		comments = append(comments, com)
	}

	if addCommentWithTrelloLink {
		cc := ch.CreateComment{
			CreatedAt: time.Now(),
			Text:      fmt.Sprintf("Card imported from Trello: %s", card.ShortURL),
		}

		comments = append(comments, cc)
	}

	return &comments
}

func buildTasks(card *Card) *[]ch.CreateTask {
	var tasks []ch.CreateTask

	for _, t := range card.Tasks {
		ts := ch.CreateTask{
			Complete:    t.Completed,
			Description: t.Description,
		}

		tasks = append(tasks, ts)
	}

	return &tasks
}

func buildLabels(card *Card) *[]ch.CreateLabel {
	var labels []ch.CreateLabel

	for _, l := range card.Labels {
		labels = append(labels, ch.CreateLabel{Name: l})
	}

	return &labels
}
