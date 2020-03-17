package quiz

import (
	"math/rand"
	"strings"
	"time"
)

type quiz struct {
	videoSource    string
	title          string
	answerVariants []string
	start          int
}

func getShuffledQuizzes(quizzes []quiz) []quiz {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]quiz, len(quizzes))
	perm := r.Perm(len(quizzes))
	for i, randIndex := range perm {
		ret[i] = quizzes[randIndex]
	}
	return ret
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
		start:          0,
	},
	{
		videoSource:    "-77UEct0cZM",
		title:          "Boku no Hero Academia (Season 1)",
		answerVariants: []string{"my hero academia", "boku no hero academia", "mha", "bnha", "my hero academy"},
		start:          0,
	},
	{
		videoSource:    "AgBUP8TJqV8",
		title:          "Attack On Titan",
		answerVariants: []string{"attack on titan", "shingeki no kyojin", "titan"},
		start:          0,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "Overlord",
		answerVariants: []string{"overlord"},
		start:          32,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "Darling in the FranXX",
		answerVariants: []string{"darling", "franxx"},
		start:          52,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "Gabriel Dropout",
		answerVariants: []string{"gabriel", "gabriel dropout"},
		start:          73,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "Code Geass",
		answerVariants: []string{"code geass"},
		start:          93,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "Re:Zero",
		answerVariants: []string{"re zero"},
		start:          113,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "Nanatsu no Taizai",
		answerVariants: []string{"nanatsu"},
		start:          134,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "No Game No Life",
		answerVariants: []string{"no game"},
		start:          151,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "Mirai Nikki",
		answerVariants: []string{"mirai nikki"},
		start:          169,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "High School of the Dead",
		answerVariants: []string{"high school"},
		start:          190,
	},
	{
		videoSource:    "Gg29TW27E2w",
		title:          "Dome x Kano",
		answerVariants: []string{"dome x kano"},
		start:          230,
	},
}
