package google

type Response struct {
	Kind    string         `json:"kind,omitempty"`
	Url     interface{}    `json:"url,omitempty"`
	Queries interface{}    `json:"queries,omitempty"`
	Context interface{}    `json:"context,omitempty"`
	Info    interface{}    `json:"searchInformation,omitempty"`
	Items   []ResponseItem `json:"items,omitempty"`
}

type ResponseItem struct {
	Kind        string      `json:"kind"`
	Title       string      `json:"title"`
	HTMLTitle   string      `json:"htmlTitle"`
	Link        string      `json:"link"`
	DisplayLink string      `json:"displayLink"`
	Snippet     string      `json:"snippet"`
	HTMLSnippet string      `json:"htmlSnippet"`
	Mime        string      `json:"string"`
	Image       interface{} `json:"image"`
}
