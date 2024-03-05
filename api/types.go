package api

type APIErrorResponse struct {
	Object  string `json:"object,omitempty"`
	Status  int    `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ResultsWrapper is the structure of your successful response
type ResultsWrapper struct {
	Results []Block `json:"results"`
}

type Block struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Heading1  *Heading   `json:"heading_1,omitempty"`
	Heading2  *Heading   `json:"heading_2,omitempty"`
	Heading3  *Heading   `json:"heading_3,omitempty"`
	Paragraph *Paragraph `json:"paragraph,omitempty"`
}

// Heading represents a generic heading, which can be used for both heading_1, heading_2, heading_3 etc.
type Heading struct {
	RichText []RichText `json:"rich_text"`
}
type Paragraph struct {
	RichText []RichText `json:"rich_text"`
}

type RichText struct {
	Type string `json:"type"`
	Text struct {
		Content string  `json:"content"`
		Link    *string `json:"link,omitempty"`
	} `json:"text"`
	Annotations struct {
		Bold          bool   `json:"bold"`
		Italic        bool   `json:"italic"`
		Strikethrough bool   `json:"strikethrough"`
		Underline     bool   `json:"underline"`
		Code          bool   `json:"code"`
		Color         string `json:"color"`
	} `json:"annotations"`
	PlainText string  `json:"plain_text"`
	Href      *string `json:"href,omitempty"`
}

// RichTextProvider interface for blocks that contain Rich Text
type RichTextProvider interface {
	GetRichText() []RichText
}
