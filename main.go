package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"web-dl/db"
	"web-dl/repository"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not load .env file")
	}

    db, err := db.GetConn()
    if err != nil {
        log.Fatal(err)
    }

    migration := repository.NewMigrationRepository(db)
    err = migration.Migrate()
    if err != nil {
        log.Fatal(err)
    }

    repo := repository.NewSourceRepository(db)
    sources, err := repo.GetSources()
    if err != nil {
        log.Fatal(err)
    }

    var wg sync.WaitGroup
	for _, source := range sources {
		wg.Add(1)

		err, links := getSources(source)
		if err != nil {
			log.Fatal(err)
		}

        source := source
		go func() {
			defer wg.Done()
			err = downloadSource(source, links)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

    wg.Wait()
}

func getSources(source *repository.Source) (error, []string) {
	res, err := http.Get(source.Url)
	if err != nil {
		return err, nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status), nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err, nil
	}

	urls := []string{}
	doc.Find(source.Selector).Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("href")
		if exists {
			urls = append(urls, val)
		}
	})

	return nil, urls
}

func downloadSource(source *repository.Source, urls []string) error {
	cmd := exec.Command("yt-dlp",
		"--throttled-rate",
		"50K",
		"-N 4",
		"-P",
		fmt.Sprintf("%s/%s", os.Getenv("DESTINATION"), source.Name),
		"-a",
		"-",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}

	for i, v := range urls {
		path := fmt.Sprintf("%s%s\n", source.Prefix, v)
		_, err = stdin.Write([]byte(path))
		if err != nil {
			return err
		}

		if i == 5 {
			break
		}
	}

	err = stdin.Close()
	if err != nil {
		return err
	}

	_ = cmd.Wait()

	return nil
}
