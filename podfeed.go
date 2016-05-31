package main

import (
	"log"
	"net/url"
	"os"
	"sort"
	"text/template"
	"time"

	"github.com/mdlayher/taggolib"
)

const rssTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>{{.FeedTitle}}</title>
    <description>{{.FeedDesc}}</description>
    {{range .Episodes}}
    <item>
      <title>{{.Title}}</title>
      <pubDate>{{.FormattedDate}}</pubDate>
      <enclosure url="{{.URL}}" length="{{.Size}}" type="audio/mpeg" />
      <guid isPermaLink="false">{{.URL}}</guid>
    </item>
    {{end}}
  </channel>
</rss>
`

var tmpl = template.Must(template.New("podfeed").Parse(rssTemplate))

type Episode struct {
	Title string
	Date  time.Time
	URL   string
	Size  int64
}

func (e Episode) FormattedDate() string {
	return e.Date.Format(time.RFC1123Z)
}

type Feed struct {
	FeedTitle string
	FeedDesc  string
	Episodes  []Episode
}

func (f Feed) Len() int           { return len(f.Episodes) }
func (f Feed) Swap(i, j int)      { f.Episodes[i], f.Episodes[j] = f.Episodes[j], f.Episodes[i] }
func (f Feed) Less(i, j int) bool { return f.Episodes[i].Date.Before(f.Episodes[j].Date) }

func main() {
	if len(os.Args) < 4 {
		log.Fatal(`Usage: podfeed "Feed Title" "Feed Description" "base URL" *.mp3`)
	}

	feed := Feed{
		FeedTitle: os.Args[1],
		FeedDesc:  os.Args[2],
	}

	baseURL, err := url.Parse(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	files := os.Args[4:]

	for _, fname := range files {
		s, err := os.Stat(fname)
		if err != nil {
			log.Fatalf("%s: %s", fname, err)
		}

		f, err := os.Open(fname)
		if err != nil {
			log.Fatalf("%s: %s", fname, err)
		}

		var title string
		tp, err := taggolib.New(f)
		if err != nil {
			log.Printf("%s: Error loading tags: %s.  Falling back to filename for title", fname, err)
			title = fname
		} else {
			title = tp.Title()
		}

		u := *baseURL
		u.Path += "/" + fname

		e := Episode{
			Title: title,
			Date:  s.ModTime(),
			URL:   u.String(),
			Size:  s.Size(),
		}

		feed.Episodes = append(feed.Episodes, e)
		f.Close()
	}

	sort.Sort(feed)
	tmpl.Execute(os.Stdout, feed)
}
