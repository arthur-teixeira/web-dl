# Web-dl
This project was made because YouTube removed the possibility of using ad blocks, and I hate ads.
The goal is to download the most recent videos of channels that I follow, and serve them in a Plex server to my local network.

## Dependencies
- [yt-dlp](https://github.com/yt-dlp/yt-dlp). There is a binary available in the dist/ folder.
- [ffmpeg](https://ffmpeg.org/)
- [Docker](https://www.docker.com/)
- [Selenium](https://www.selenium.dev/). There is a script in the dist/ folder that downloads the correct version of all the drivers needed. (Taken from the [go selenium package](https://github.com/tebeka/selenium/blob/master/vendor/init.go)).
- [Optional] [Plex server](https://www.plex.tv/pt-br/media-server-downloads/).

## Quick start

```
$ cd postgres
$ docker compose up -d
$ cd ../
$ go run main.go
```

## Notes
- Currently, you need to add video sources manually to the database. In the future, I plan on making a web UI to add and manage them.
- You can automatically download videos from your sources by building a binary with go build and running it as a Cron job. In the future, I plan on adding a Cron scheduler to the project itself.
