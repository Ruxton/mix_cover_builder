package google

type Response struct {
	Data    ResponseData `json:"responseData,omitempty"`
	Details interface{}  `json:"responseDetails,omitempty"`
	Status  int          `json:"responseStatus,omitempty"`
}

type ResponseData struct {
	Results []Result   `json:"results,omitempty"`
	Cursor  CursorData `json:"cursor,omitempty"`
}

type CursorData struct {
	Count            string `json:"resultCount,omitempty"`
	Pages            []Page `json:"pages,omitempty"`
	EstimatedResults string `json:estimatedResultCount,omitempty`
	CurrentPage      int    `json:currentPageIndex,omitempty`
	MoreResultsURL   string `json:moreResultsUrl,omitempty`
	SearchResultTime string `json:searchResultTime,omitempty`
}

type Page struct {
	Start string `json:"start,omitempty"`
	Label int    `json:"label,omitempty"`
}

type Result struct {
	GsearchResultClass  string
	Width               string
	Height              string
	ImageId             string
	TbWidth             string
	TbHeight            string
	UnescapedUrl        string
	Url                 string
	VisibleUrl          string
	Title               string
	TitleNoFormatting   string
	OriginalContextUrl  string
	Content             string
	ContentNoFormatting string
	TbUrl               string
}
