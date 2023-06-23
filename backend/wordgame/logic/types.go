package logic

type WordGameEntry struct {
	CorrectWord         string   `json:"correctWord"`
	CorrectWordImageUrl string   `json:"correctWordImageUrl"`
	IncorrectWords      []string `json:"incorrectWords"`
}

type WordGameSubmission struct {
	Username        string `json:"username"`
	IsAnswerCorrect bool   `json:"isAnswerCorrect"`
	Word            string `json:"word"`
}

type WordGameUserInfo struct {
	Username string `json:"username"`
}
