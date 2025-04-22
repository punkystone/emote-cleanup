package internal

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type SortedEmote struct {
	ID       string
	Name     string
	Count    int
	Added    time.Time
	LastUsed *time.Time
	Score    float64
}

const hoursInDay = 24

func Render(emotesCount map[string]*Emote, renderFile string) error {
	sortedEmotes := sortEmotes(emotesCount)

	htmlBuilder := strings.Builder{}
	htmlBuilder.WriteString(htmlHeader)
	htmlBuilder.WriteString(fmt.Sprintf("<h4>%s</h4>", time.Now().Format("02.01.2006")))
	htmlBuilder.WriteString(tableHeader)
	for _, emote := range sortedEmotes {
		htmlBuilder.WriteString(fmt.Sprintf(`<tr>
		<td><div><img src='https://cdn.7tv.app/emote/%s/1x.avif' loading="lazy"></div></td>
		<td><div>%s</td>
		<td>%d</td>
		<td>%d Tagen</td>
		<td>%s</td>
		<td>%f</td>
		<td><div><a href='https://7tv.app/emotes/%s'>%s</div></a></td>
		</tr>`,
			emote.ID, emote.Name, emote.Count, int(time.Since(emote.Added).Hours()/hoursInDay), formatLastUsed(emote.LastUsed), emote.Score, emote.ID, openIcon))
	}
	htmlBuilder.WriteString(htmlFooter)
	html := htmlBuilder.String()
	html = strings.ReplaceAll(html, "\n", "")
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
		if len(emote.LastUsed) == 0 {
			sortedEmotes = append(sortedEmotes, SortedEmote{
				ID:       emote.ID,
				Name:     emoteName,
				Added:    emote.Added,
				Count:    0,
				LastUsed: nil,
				Score:    emote.Score,
			})
			continue
		}

		sort.Slice(emote.LastUsed, func(i, j int) bool {
			return emote.LastUsed[i].Before(*emote.LastUsed[j])
		})
		lastUsed := emote.LastUsed[len(emote.LastUsed)-1]
		sortedEmotes = append(sortedEmotes, SortedEmote{
			ID:       emote.ID,
			Name:     emoteName,
			Added:    emote.Added,
			Count:    len(emote.LastUsed),
			LastUsed: lastUsed,
			Score:    emote.Score,
		})
	}
	sort.Slice(sortedEmotes, func(i, j int) bool {
		return sortedEmotes[i].Score > sortedEmotes[j].Score
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
		flex-direction: column;
	}
	.main > h2 {
		margin-bottom: 20px;
	}
	.main > h4 {
		margin-bottom: 20px;
	}
    table {
        border-collapse: collapse;
        box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
		table-layout: auto !important;
    }
    table th, td {
        border: 1px solid black;
        padding: 10px;
    }
	img {
		height: 30px;
	}
	td:nth-child(1) > div {
	 	display: flex;
	 	align-items: center;
		justify-content: center;
	}
	td:nth-child(2) > div {
		display: flex;
        align-items: center;
	}
	td:nth-child(2) > div > a {
		padding-left: 10px;
	}

    
    tr:nth-of-type(even) {
        background-color: #c6c6c6;
    }

	 td:nth-child(3) {
        text-align: center;
    }
    
    </style>
    </head>
    <body>
    <div class="main">
	<h2>Emote Scoreboard</h2>
`
const tableHeader = `
	<table>
    <tr>
	<th>Emote</th>
    <th>Name</th>
    <th>Count</th>
    <th>Hinzugefügt vor</th>
    <th>Zuletzt benutzt</th>
	<th>Score</th>
    </tr>
`

const htmlFooter = `
	</table>
    </div>
    </body>
    </html>
`
const openIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#000" d="M19 19H5V5h7V3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14c1.1 0 2-.9 2-2v-7h-2zM14 3v2h3.59l-9.83 9.83l1.41 1.41L19 6.41V10h2V3z"/></svg>`
