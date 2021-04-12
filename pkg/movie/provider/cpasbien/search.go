package cpasbien

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/TheoBrigitte/confluence/pkg/movie"
	anatorrent "github.com/anacrolix/torrent"
	"golang.org/x/sync/errgroup"
)

func (c *client) Search(query string) ([]movie.MovieTorrent, error) {
	escapedQuery := url.PathEscape(query)
	searchURL := fmt.Sprintf(searchURLFormat, escapedQuery)
	u, err := baseURL.Parse(searchURL)
	if err != nil {
		return nil, err
	}
	log.Printf("request %s %#q", "GET", u.String())
	res, err := c.http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	log.Printf("response %s %#q body=%d %v", res.Status, u.String(), res.ContentLength, res.Header)

	links, err := c.searchProcess(res.Body)
	if err != nil {
		return nil, err
	}
	if len(links) > 10 {
		links = links[:10]
	}

	return c.getDetails(links)
}

func (c *client) searchProcess(r io.ReadCloser) ([]*url.URL, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	return c.getLinks(doc)
}

func (c *client) getLinks(doc *goquery.Document) (links []*url.URL, mainErr error) {
	linkSelector := "#conteneur #gauche div.table_div table.table-corps tbody tr td a.titre"
	doc.Find(linkSelector).EachWithBreak(func(i int, s *goquery.Selection) bool {
		val, exist := s.Attr("href")
		if !exist {
			mainErr = fmt.Errorf("getLinks: no link: %v", s)
			return false
		}

		u, err := url.Parse(val)
		if err != nil {
			mainErr = err
			return false
		}

		links = append(links, baseURL.ResolveReference(u))

		return true
	})

	return links, nil
}

func (c *client) getDetails(urls []*url.URL) (movies []movie.MovieTorrent, err error) {
	var (
		g         errgroup.Group
		movieChan = make(chan movie.MovieTorrent)
	)

	go func() {
		for i := 0; i < len(urls); i++ {
			movies = append(movies, <-movieChan)
		}
	}()
	for _, url := range urls {
		u := url
		g.Go(func() error {
			return c.movieDetail(u, movieChan)
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (c *client) movieDetail(url *url.URL, movieChan chan movie.MovieTorrent) error {
	movie := &movie.MovieTorrent{}
	defer func() { movieChan <- *movie }()

	log.Printf("request %s %#q", "GET", url.String())
	res, err := c.http.Get(url.String())
	log.Printf("response %s %#q body=%d %v", res.Status, url.String(), res.ContentLength, res.Header)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	m, err := c.movieProcess(res.Body)
	if err != nil {
		return err
	}
	movie = m
	movie.Torrent.URL = url.String()

	parts := strings.Split(strings.TrimSuffix(url.String(), ".html"), "-")
	id := parts[len(parts)-1]
	movie.MovieBase.ID, err = strconv.Atoi(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) movieProcess(r io.ReadCloser) (*movie.MovieTorrent, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	hash, err := movieHash(doc)
	if err != nil {
		return nil, err
	}
	m := &movie.MovieTorrent{
		MovieBase: movie.MovieBase{
			Title: movieTitle(doc),
		},
		Torrent: movie.Torrent{
			Hash: hash,
		},
	}

	return m, nil
}

func movieTitle(doc *goquery.Document) string {
	titleSelector := "#gauche > div.h2fiche > a > h1"
	return strings.TrimSpace(doc.Find(titleSelector).Text())
}

func movieHash(doc *goquery.Document) (string, error) {
	hashSelector := "#textefiche > p:nth-child(2) > a:nth-child(2)"
	val, _ := doc.Find(hashSelector).Attr("href")
	t, err := anatorrent.TorrentSpecFromMagnetURI(val)
	if err != nil {
		return "", err
	}
	return t.InfoHash.String(), nil
}
