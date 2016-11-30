package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	trello "github.com/jnormington/go-trello"
	dropbox "github.com/tj/go-dropbox"
)

var dateLayout = "2006-01-02T15:04:05.000Z"
var safeFileNameRegexp = regexp.MustCompile(`[^a-zA-Z0-9_.]+`)

// Card holds all the attributes needed for migrating a complete card from Trello to Clubhouse
type Card struct {
	Name        string            `json:"name"`
	Desc        string            `json:"desc"`
	Labels      []string          `json:"labels"`
	DueDate     *time.Time        `json:"due_date"`
	IdCreator   string            `json:"id_creator"`
	CreatedAt   *time.Time        `json:"created_at"`
	Comments    []Comment         `json:"comments"`
	Tasks       []Task            `json:"checklists"`
	Position    float32           `json:"position"`
	ShortURL    string            `json:"url"`
	Attachments map[string]string `json:"attachments"`
	IdMembers   []string          `json:"id_members"`
}

// Task builds a basic object based off trello.Task
type Task struct {
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
}

// Comment builds a basic object based off trello.Comment
type Comment struct {
	Text      string
	Creator   string
	CreatedAt *time.Time
}

// ProcessCardsForExporting takes *[]trello.Card, *TrelloOptions and builds up a Card
// which consists of calling other functions to make the api calls to Trello
// for the relevant attributes of a card returns *[]Card
func ProcessCardsForExporting(crds *[]trello.Card, opts *TrelloOptions) *[]Card {
	var cards []Card

	for _, card := range *crds {
		var c Card

		c.Name = card.Name
		c.Desc = card.Desc
		c.Labels = getLabelsFlattenFromCard(&card)
		c.DueDate = parseDateOrReturnNil(card.Due)
		c.Creator, c.CreatedAt, c.Comments = getCommentsAndCardCreator(&card)
		c.IdMembers = card.IdMembers
		c.Tasks = getCheckListsForCard(&card)
		c.Position = card.Pos
		c.ShortURL = card.ShortUrl

		if opts.ProcessImages {
			c.Attachments = downloadCardAttachmentsUploadToDropbox(&card)
		}

		cards = append(cards, c)
	}

	return &cards
}

func getCommentsAndCardCreator(card *trello.Card) (string, *time.Time, []Comment) {
	var creator string
	var createdAt *time.Time
	var comments []Comment

	actions, err := card.Actions()
	if err != nil {
		fmt.Println("Error: Querying the actions for:", card.Name, "ignoring...", err)
	}

	for _, a := range actions {
		if a.Type == "commentCard" && a.Data.Text != "" {
			c := Comment{
				Text:      a.Data.Text,
				Creator:   a.IdMemberCreator,
				CreatedAt: parseDateOrReturnNil(a.Date),
			}
			comments = append(comments, c)

		} else if a.Type == "createCard" {
			creator = a.MemberCreator.FullName
			createdAt = parseDateOrReturnNil(a.Date)
		}
	}

	return creator, createdAt, comments
}

func getCheckListsForCard(card *trello.Card) []Task {
	var tasks []Task

	checklists, err := card.Checklists()
	if err != nil {
		fmt.Println("Error: Occurred querying checklists for:", card.Name, "ignoring...", err)
	}

	for _, cl := range checklists {
		for _, i := range cl.CheckItems {
			var completed bool
			if i.State == "complete" {
				completed = true
			}

			t := Task{
				Completed:   completed,
				Description: fmt.Sprintf("%s - %s", cl.Name, i.Name),
			}

			tasks = append(tasks, t)
		}
	}

	return tasks
}

func getLabelsFlattenFromCard(card *trello.Card) []string {
	var labels []string

	for _, l := range card.Labels {
		labels = append(labels, l.Name)
	}

	return labels
}

func parseDateOrReturnNil(strDate string) *time.Time {
	d, err := time.Parse(dateLayout, strDate)
	if err != nil {
		//If the date isn't parseable from trello api just return nil
		return nil
	}

	return &d
}

func downloadCardAttachmentsUploadToDropbox(card *trello.Card) map[string]string {
	sharedLinks := map[string]string{}
	d := dropbox.New(dropbox.NewConfig(dropboxToken))

	attachments, err := card.Attachments()
	if err != nil {
		log.Fatal(err)
	}

	for i, f := range attachments {
		name := safeFileNameRegexp.ReplaceAllString(f.Name, "_")
		path := fmt.Sprintf("/trello/%s/%s/%d%s%s", card.IdList, card.Id, i, "_", name)

		io := downloadTrelloAttachment(&f)
		_, err := d.Files.Upload(&dropbox.UploadInput{
			Path:   path,
			Mode:   dropbox.WriteModeAdd,
			Reader: io,
			Mute:   true,
		})

		io.Close()

		if err != nil {
			fmt.Printf("Error occurred uploading file to dropbox continuing... %s\n", err)
		} else {
			// Must be success created a shared url
			s := dropbox.CreateSharedLinkInput{path, false}
			out, err := d.Sharing.CreateSharedLink(&s)
			if err != nil {
				fmt.Printf("Error occurred sharing file on dropbox continuing... %s\n", err)
			} else {
				sharedLinks[name] = out.URL
			}
		}
	}

	return sharedLinks
}

func downloadTrelloAttachment(attachment *trello.Attachment) io.ReadCloser {
	resp, err := http.Get(attachment.Url)
	//	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("Error in download Trello attachment %s\n", err)
	}

	return resp.Body
}
