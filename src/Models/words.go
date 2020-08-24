package Models

type AddWordRequest struct {
	Word string `json:"word"`
}

type AddWordResponse struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	CurrentSentence string `json:"current_sentence"`
}

