package logic

type WordGameEntry struct {
	CorrectWord         string   `json:"correctWord"`
	CorrectWordImageUrl string   `json:"correctWordImageUrl"`
	IncorrectWords      []string `json:"incorrectWords"`
}

type WordGameSubmission struct {
	IsAnswerCorrect bool   `json:"isAnswerCorrect"`
	Word            string `json:"word"`
}
