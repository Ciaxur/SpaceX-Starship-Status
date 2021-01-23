package main

// TweetData Data Key
type TweetData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// TweetMeta Request Meta Data
type TweetMeta struct {
	OldestID    string `json:"oldest_id"`
	NewestID    string `json:"newest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// Tweet Reply Interface
type Tweet struct {
	Data        []TweetData `json:"data"`         // Tweet Data
	Meta        TweetMeta   `json:"meta"`         // Metadata
	LatestMatch TweetData   `json:"latest_match"` // Object of the Latest Match
}
