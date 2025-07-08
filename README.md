# YouTube CLI (`cliyt`)

A minimalistic and fast command-line interface (CLI) tool to search and play YouTube videos using the official YouTube Data API and `mpv` media player.

Built in Go — Simple. Clean. Efficient.

---

## Features

* Search YouTube videos directly from your terminal
* Stream videos with `mpv` in one command
* Displays video duration, title, and channel name
* Customizable result limit (`setMax`)
* Terminal screen clearing (`clear`)
* Minimal dependencies, fast execution

---

## Getting Started

### Prerequisites

* [Go](https://go.dev/) installed (Go 1.20+ recommended)
* [mpv](https://mpv.io/) media player installed and accessible in `$PATH`
* A valid YouTube Data API v3 Key

---

### Installation

```bash
git clone https://github.com/yourusername/cliyt.git
cd cliyt
go build -o cliyt
./cliyt
```

Replace the hardcoded API key in `apiRequest()` and `getVideoDuration()` with your own.

---

## Usage

Upon running the app, you'll see a prompt:

```
Welcome to YouTube CLI!
Commands:
  search <query> - Search for videos
  exit           - Exit the program
```

### Search for videos

```bash
cliyt > search lo-fi beats
```

Displays top 5 YouTube results (title, channel name, and duration).

### Play a video

```bash
cliyt > 2
```

Plays the 2nd video using `mpv`.

### Set maximum results

```bash
cliyt > setMax 10
```

Sets max results for future searches (1–50).

### Clear screen

```bash
cliyt > clear
```

Clears the terminal output.

### Exit the CLI

```bash
cliyt > exit
```

---

## Security Notice

This tool hardcodes an API key, which is insecure for production use. To secure your key:

* Load it from an environment variable or `.env` file
* Avoid committing your key to version control

---

## Design Highlights

* Uses YouTube Data API v3 for fetching search results and durations
* Parses ISO 8601 duration format (e.g., `PT1H2M3S`) into `hh:mm:ss`
* Supports color-coded terminal output using ANSI escape sequences
* Pipes video URLs to `mpv` for seamless playback

---

## Built With

* [Go](https://go.dev/) – Fast and simple compiled language
* [YouTube Data API v3](https://developers.google.com/youtube/v3)
* [mpv](https://mpv.io/) – Cross-platform media player

---

## TODO

* [ ] Add support for paginated results
* [ ] Add option to download videos using `yt-dlp`
* [ ] Secure API key handling
* [ ] Support for playlists or trending videos

---

## License

MIT License. Use it, modify it, ship it.

---

## Contributing

Pull requests are welcome. For major changes, open an issue first to discuss what you’d like to change.

---

## Example Session

```
Welcome to YouTube CLI!
Commands:
  search <query> - Search for videos
  exit           - Exit the program

cliyt > search lo-fi beats

1. lo-fi chill beats to relax/study to
   Channel: ChilledCow
   Duration: 1:00:00
----------------------------------
2. lo-fi hip hop radio - beats to sleep/chill to
   Channel: Chillhop Music
   Duration: 2:30:00
----------------------------------

cliyt > 1
Now playing, lo-fi chill beats to relax/study to
```

---

Let me know if you'd like to modularize the codebase, support more players, or add `.env` configuration.

