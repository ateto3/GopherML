package chat

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Rule struct{ Re *regexp.Regexp; Replies []string }
type Bot struct{ Rules []Rule; rand *rand.Rand }

func New(rules []Rule) *Bot { return &Bot{Rules: rules, rand: rand.New(rand.NewSource(time.Now().UnixNano()))} }

func (b *Bot) Reply(input string) string {
	in := strings.TrimSpace(strings.ToLower(input))
	for _, r := range b.Rules {
		if m := r.Re.FindStringSubmatch(in); m != nil {
			resp := r.Replies[b.rand.Intn(len(r.Replies))]
			for i := 1; i < len(m); i++ { resp = strings.ReplaceAll(resp, "%"+itoa(i), reflect(m[i])) }
			return resp
		}
	}
	return "Tell me more."
}

func reflect(s string) string {
	r := strings.NewReplacer(" i ", " you ", " me ", " you ", " my ", " your ", " am ", " are ")
	return r.Replace(" " + s + " ")[1:len(s)+1]
}

func itoa(i int) string { return string(rune('0' + i)) }

func Eliza() *Bot {
	return New([]Rule{
		{regexp.MustCompile(`^i need (.*)`), []string{"Why do you need %1?", "Would it really help to get %1?"}},
		{regexp.MustCompile(`^i am (.*)`), []string{"How long have you been %1?", "How do you feel about being %1?"}},
		{regexp.MustCompile(`^i can't (.*)`), []string{"What would it take for you to %1?", "Perhaps you could %1 if you tried."}},
		{regexp.MustCompile(`^why (.*)`), []string{"Why do you think %1?", "What do you think?"}},
		{regexp.MustCompile(`^hello(.*)`), []string{"Hello. How are you?", "Hi there."}},
	})
}
