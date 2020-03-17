package quiz

import "strings"

type quiz struct {
	videoSource    string
	title          string
	answerVariants []string
}

func (q *quiz) isAnswerRight(answer string) bool {
	answer = strings.ToLower(strings.Trim(answer, " ")) // todo: write test cases
	for _, v := range q.answerVariants {
		if v == answer {
			return true
		}
	}
	return false
}

var hardcodedQuizzes = []quiz{
	{
		videoSource:    "t-QSmNReDyI",
		title:          "Neon Genesis Evangelion",
		answerVariants: []string{"neon genesis evangelion", "evangelion", "eva", "nge"},
	},
	{
		videoSource:    "-77UEct0cZM",
		title:          "Boku no Hero Academia (Season 1)",
		answerVariants: []string{"my hero academia", "boku no hero academia", "mha", "bnha", "my hero academy"},
	},
	{
		videoSource:    "AgBUP8TJqV8",
		title:          "Attack On Titan",
		answerVariants: []string{"attack on titan", "shingeki no kyojin", "titan", "snk"},
	},
}
