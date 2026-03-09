package reports

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
)

type notionClient struct {
	client *notionapi.Client
}

func NewNotionClient(notionToken string) NotionClient {
	client := notionapi.NewClient(notionapi.Token(notionToken))
	return &notionClient{client: client}
}

func (n *notionClient) GetOrCreateReportsPage(ctx context.Context, dbID notionapi.DatabaseID) (notionapi.PageID, error) {
	// Create page
	page, err := n.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: dbID,
		},
		Properties: notionapi.Properties{
			"Candidate": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: "Interview of Sourabh Mandal",
						},
					},
				},
			},
			"Interview Date": notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: (*notionapi.Date)(new(notionapi.Date)),
				},
			},
			"Interviewer": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: "Interview of Sourabh Mandal",
						},
					},
				},
			},
			"Score": notionapi.NumberProperty{
				Number: 8,
			},
		},
	})

	if err != nil {
		return "", err
	}

	return notionapi.PageID(page.ID), nil
}

func (n *notionClient) GetOrCreateReportsDatabase(ctx context.Context, pageID notionapi.PageID) (notionapi.DatabaseID, error) {
	searchResp, err := n.client.Search.Do(ctx, &notionapi.SearchRequest{
		Query: "Reports",
		Filter: notionapi.SearchFilter{
			Value:    "data_source",
			Property: "object",
		},
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(searchResp.Results)

	for _, res := range searchResp.Results {
		if db, ok := res.(*notionapi.Database); ok {
			return notionapi.DatabaseID(db.ID), nil
		}
	}

	db, err := n.client.Database.Create(ctx, &notionapi.DatabaseCreateRequest{
		IsInline: true,
		Parent: notionapi.Parent{
			PageID: pageID,
		},
		Title: []notionapi.RichText{
			{
				Text: &notionapi.Text{
					Content: "Reports",
				},
			},
		},
		Properties: notionapi.PropertyConfigs{
			"Candidate": notionapi.TitlePropertyConfig{
				ID:    "candidate",
				Type:  "title",
				Title: struct{}{},
			},
			"Score": notionapi.NumberPropertyConfig{
				ID:   "score",
				Type: "number",
				Number: notionapi.NumberFormat{
					Format: "number",
				},
			},
			"Interview Date": notionapi.DatePropertyConfig{
				ID:   "interview_date",
				Type: "date",
				Date: struct{}{},
			},
			"Interviewer": notionapi.RichTextPropertyConfig{
				ID:       "interviewer",
				Type:     "rich_text",
				RichText: struct{}{},
			},
		},
	})

	if err != nil {
		return "", err
	}

	fmt.Println(db)

	return notionapi.DatabaseID(db.ID), nil
}
