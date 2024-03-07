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
	Todo      *Todo      `json:"to_do,omitempty"`
	Bookmark  *Bookmark  `json:"bookmark,omitempty"`
	Bulleted  *ListItem  `json:"bulleted_list_item,omitempty"`
	Numbered  *ListItem  `json:"numbered_list_item,omitempty"`
	Paragraph *Paragraph `json:"paragraph,omitempty"`
	Code      *Code      `json:"code,omitempty"`
}

// Heading represents a generic heading, which can be used for both heading_1, heading_2, heading_3 etc.
type Heading struct {
	RichText []RichText `json:"rich_text"`
}
type Paragraph struct {
	RichText []RichText `json:"rich_text"`
}

type ListItem struct {
	RichText []RichText `json:"rich_text"`
}

type Bookmark struct {
	URL string `json:"url"`
}

type Code struct {
	RichText []RichText `json:"rich_text"`
	Language string     `json:"language"`
}

type Todo struct {
	RichText []RichText `json:"rich_text"`
	Checked  bool       `json:"checked"`
}

type LinkObject struct {
	URL *string `json:"url,omitempty"`
}
type Text struct {
	Content string      `json:"content"`
	Link    *LinkObject `json:"link,omitempty"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type RichText struct {
	Type        string      `json:"type"`
	Text        Text        `json:"text"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        *string     `json:"href,omitempty"`
}

// RichTextProvider interface for blocks that contain Rich Text
type RichTextProvider interface {
	GetRichText() []RichText
}

// Implement GetRichText for Heading
func (h *Heading) GetRichText() []RichText {
	return h.RichText
}

// Implement GetRichText for Paragraph
func (p *Paragraph) GetRichText() []RichText {
	return p.RichText
}

// Implement GetRichText for ListItem
func (p *ListItem) GetRichText() []RichText {
	return p.RichText
}

// Implement GetRichText for to_do
func (t *Todo) GetRichText() []RichText {
	return t.RichText
}

// Implement GetRichText for code
func (c *Code) GetRichText() []RichText {
	return c.RichText
}
