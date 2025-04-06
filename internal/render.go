package internal

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type SortedEmote struct {
	Name     string
	Count    int
	Added    time.Time
	LastUsed *time.Time
}

const hoursInDay = 24

func Render(emotesCount map[string]*Emote, renderFile string) error {
	sortedEmotes := sortEmotes(emotesCount)

	htmlBuilder := strings.Builder{}
	htmlBuilder.WriteString(htmlHeader)
	for _, emote := range sortedEmotes {
		htmlBuilder.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%d Tagen</td><td>%s</td></tr>", emote.Name, emote.Count, int(time.Since(emote.Added).Hours()/hoursInDay), formatLastUsed(emote.LastUsed)))
	}
	htmlBuilder.WriteString(htmlFooter)
	html := htmlBuilder.String()
	err := writeFile(renderFile, html)
	if err != nil {
		return err
	}
	return nil
}

func formatLastUsed(lastUsed *time.Time) string {
	if lastUsed == nil {
		return "Nie benutzt"
	}
	days := int(time.Since(*lastUsed).Hours() / hoursInDay)
	if days == 0 {
		return "Heute benutzt"
	}
	return fmt.Sprintf("%d Tagen", days)
}

func writeFile(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

func sortEmotes(emotes map[string]*Emote) []SortedEmote {
	sortedEmotes := make([]SortedEmote, 0, len(emotes))
	for emoteName, emote := range emotes {
		sortedEmotes = append(sortedEmotes, SortedEmote{
			Name:     emoteName,
			Count:    emote.Count,
			Added:    emote.Added,
			LastUsed: emote.LastUsed,
		})
	}
	sort.Slice(sortedEmotes, func(i, j int) bool {
		return sortedEmotes[i].Count > sortedEmotes[j].Count
	})
	return sortedEmotes
}

const htmlHeader = `
<html>
    <head>
    <title>Emote Usage</title>
    <style>
    * {
        font-family: Arial, Helvetica, sans-serif;
        padding: 0;
        margin: 0;
        box-sizing: border-box;
    }
    .main {
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 30px;
        }
    table {
        border-collapse: collapse;
        min-width: 50%;
        box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
        }
       table th, td {
        border: 1px solid black;
        padding: 10px;
    }

    
    tr:nth-of-type(even) {
        background-color: #c6c6c6;
    }

    td:nth-child(2), td:nth-child(3) {
        text-align: center;
    }
    
    </style>
    </head>
    <body>
    <div class="main">
    <table>
    <tr>
    <th>Emote</th>
    <th>Count</th>
    <th>Hinzugefügt vor</th>
    <th>Zuletzt benutzt</th>
    </tr>
`

const htmlFooter = `
	</table>
    </div>
    </body>
    </html>
`
