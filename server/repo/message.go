package repo

type Message struct {
	EleIdArr []string `json:"eleIdArr"`
	Url      string   `json:"url"`
}

type Messages []Message
