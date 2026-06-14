package philpapers

// Paper is a philosophy paper record from PhilPapers via OAI-PMH.
type Paper struct {
	Rank    int    `json:"rank"`
	ID      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Subject string `json:"subject"`
	URL     string `json:"url"`
}

// Category is a philosophy category from PhilPapers.
type Category struct {
	Rank int    `json:"rank"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// DefaultCategories lists the top-level philosophy categories on PhilPapers.
var DefaultCategories = []Category{
	{1, "Epistemology", "epistemology"},
	{2, "Ethics", "ethics"},
	{3, "Philosophy of Mind", "philosophy-of-mind"},
	{4, "Logic and Philosophy of Logic", "logic-and-philosophy-of-logic"},
	{5, "Metaphysics", "metaphysics"},
	{6, "Philosophy of Science", "philosophy-of-science"},
}

// Categories returns the hardcoded list of PhilPapers philosophy categories.
func Categories() []Category {
	return DefaultCategories
}

// OAI-PMH wire types for XML unmarshalling.

type oaiPMH struct {
	ListRecords *listRecords `xml:"ListRecords"`
	Error       *oaiError    `xml:"error"`
}

type listRecords struct {
	Records         []oaiRecord `xml:"record"`
	ResumptionToken string      `xml:"resumptionToken"`
}

type oaiRecord struct {
	Header   oaiHeader   `xml:"header"`
	Metadata oaiMetadata `xml:"metadata"`
}

type oaiHeader struct {
	Identifier string `xml:"identifier"`
	Datestamp  string `xml:"datestamp"`
}

type oaiMetadata struct {
	DC oaiDC `xml:"dc"`
}

type oaiDC struct {
	Titles    []string `xml:"title"`
	Creators  []string `xml:"creator"`
	Dates     []string `xml:"date"`
	Types     []string `xml:"type"`
	Relations []string `xml:"relation"`
	Subjects  []string `xml:"subject"`
}

type oaiError struct {
	Code    string `xml:"code,attr"`
	Message string `xml:",chardata"`
}
