package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"time"

	trello "github.com/jnormington/go-trello"
	dropbox "github.com/tj/go-dropbox"
)

var dateLayout = "2006-01-02T15:04:05.000Z"
var safeFileNameRegexp = regexp.MustCompile(`[^a-zA-Z0-9_.]+`)

var mimeTypes = map[string]string{
	".gif":  "image/gif",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
}

// Card holds all the attributes needed for migrating a complete card from Trello to Clubhouse
type Card struct {
	Name        string       `json:"name"`
	Desc        string       `json:"desc"`
	Labels      []string     `json:"labels"`
	DueDate     *time.Time   `json:"due_date"`
	IDCreator   string       `json:"id_creator"`
	IDOwners    []string     `json:"id_owners"`
	CreatedAt   *time.Time   `json:"created_at"`
	Comments    []Comment    `json:"comments"`
	Tasks       []Task       `json:"checklists"`
	Position    float32      `json:"position"`
	ShortURL    string       `json:"url"`
	Attachments []SharedItem `json:"attachments"`
}

// Task builds a basic object based off trello.Task
type Task struct {
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
}

// Comment builds a basic object based off trello.Comment
type Comment struct {
	Text        string
	IDCreator   string
	CreatorName string
	CreatedAt   *time.Time
}

//SharedItem represents an attachment from Trello
type SharedItem struct {
	Name     string
	URL      string
	MimeType string
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
		c.IDCreator, c.CreatedAt, c.Comments = getCommentsAndCardCreator(&card)
		c.Tasks = getCheckListsForCard(&card)
		c.Position = card.Pos
		c.ShortURL = card.ShortUrl
		c.IDOwners = card.IdMembers

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
				Text:        a.Data.Text,
				IDCreator:   a.MemberCreator.Id,
				CreatorName: a.MemberCreator.FullName,
				CreatedAt:   parseDateOrReturnNil(a.Date),
			}
			comments = append(comments, c)

		} else if a.Type == "createCard" {
			creator = a.MemberCreator.Id
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

func downloadCardAttachmentsUploadToDropbox(card *trello.Card) []SharedItem {
	var sharedItems []SharedItem
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
				u, err := url.Parse(out.URL)

				si := SharedItem{
					Name:     name,
					URL:      out.URL,
					MimeType: determineMimeType(out.Path),
				}

				if err == nil {
					// So that we can get the image inline on the ticket
					// parse it and amend the dl param to 1
					q := u.Query()
					q.Set("dl", "1")
					u.RawQuery = q.Encode()

					fmt.Println("Got here new url: %s", u.String())
					si.URL = u.String()
				}

				sharedItems = append(sharedItems, si)
				fmt.Printf("%#v\n", si)
			}
		}
	}

	return sharedItems
}

// determineMimeType supports just determine the mimeType from extention for only the most common
// image types so that a thumbnail maybe displayed inline on the cluhbouse ticket
func determineMimeType(path string) string {
	ext := filepath.Ext(path)

	return mimeTypes[ext]
}

func downloadTrelloAttachment(attachment *trello.Attachment) io.ReadCloser {
	resp, err := http.Get(attachment.Url)
	//	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("Error in download Trello attachment %s\n", err)
	}

	return resp.Body
}
