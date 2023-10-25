package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"web-dl/db"
	"web-dl/repository"

	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
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

	numWorkers := 3
	linksChan := make(chan []string, len(sources))
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = downloadSource(linksChan)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	for _, source := range sources {
		err, links := getSources(source)
		if err != nil {
			log.Fatal(err)
		}
		linksChan <- links
	}

	close(linksChan)
	wg.Wait()
}

func getSources(source *repository.Source) (error, []string) {
	const (
		seleniumPath    = "dist/selenium-server.jar"
		geckoDriverPath = "dist/geckodriver"
		port            = 8080
	)

	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}

	log.Println("opening selenium")
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		return err, nil
	}
	defer service.Stop()

	capabilities := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		return err, nil
	}
	defer wd.Quit()

	if err := wd.Get(source.Url); err != nil {
		panic(err)
	}

	elems, err := wd.FindElements(selenium.ByCSSSelector, source.Selector)
	if err != nil {
		panic(err)
	}

    var urls []string
	for _, e := range elems {
		val, err := e.GetAttribute("href")
		if err == nil {
			path := fmt.Sprintf("%s%s\n", source.Prefix, val)
			urls = append(urls, path)
		}
	}

    return nil, urls
}

func downloadSource(urls <-chan []string) error {
	cmd := exec.Command("yt-dlp",
		"--throttled-rate",
		"50K",
		"-N 4",
		"-P",
		os.Getenv("DESTINATION"),
		"-o%(uploader)s/%(title)s.%(ext)s",
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

	for i, v := range <-urls {
		_, err = stdin.Write([]byte(v))
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
