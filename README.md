# 🇫🇷 Mot du jour API

Mot du jour is a simple proxy API for to retrieve a random French word from [Wiktionary](https://en.wiktionary.org/wiki/Wiktionary:Random_page).

It returns the word in a structured JSON format, making it easier for other applications to consume (since Wiktionary itself doesn't currently offer an official API).

## 🏃 To run

Clone the repo and run

```bash
go run .
```

## ⚡ Usage

Simply send a GET request to the API's `/mot_du_jour` endpoint like so:

```bash
curl http://localhost:8080/mot_du_jour
```

You will receive a JSON output in the following format:

```json
{
  "Name": "œil poché",
  "Meanings": [
    {
      "Type": "Noun",
      "Gender": "masc.",
      "Definitions": ["black eye"]
    }
  ]
}
```
