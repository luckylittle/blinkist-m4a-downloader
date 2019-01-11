# Blinkist M4A Downloader

## Requirements

- Golang `1.10.2` or higher.
- Blinkist (https://www.blinkist.com) Premium account.
- Roughly 20 GB of free disk space (1,771 books as of Fri, Jan 11 2019, on average 10MB per book - but not all of them have audio).

## Configuration

Enter your username and password in `blinkist/main.go`, lines #16, #17.

## Application

- Run `go run main.go` inside `blinkist/` folder to produce `books_urls.txt`, the list of unduplicated URLs of all of the books.
- Run `go run download.go` inside `download/` folder to start downloading audio files from the above URLs.

## Technical details of the solution

1. Look for HTML tag `data-book-id` e.g.`"5c28f2fc6cee070008e7a3d7"` in each book URL.

2. Look for all HTML tags `data-chapterNo` e.g.`"1"` and corresponding `data-chapterId` e.g.`"5c28f3296cee070007b46369"` (both on the same line) from each book URL.

3. Construct this API link to get the short-lived download link: `https://www.blinkist.com/api/books/<data-book-id>/chapters/<data-chapterId>/audio`.
(e.g.`https://www.blinkist.com/api/books/5c28f2fc6cee070008e7a3d7/chapters/5c28f3296cee070007b46369/audio`).

4. Read the output for each book chapter, e.g.:
```json
{"url":"https://abcdefgh12345.cloudfront.net/5c28f2fc6cee070008e7a3d7/5c28f3296cee070007b46369.m4a?Expires=1234567890\u0026Signature=abcdefghijklmnopqrstuvwxyz1234-567890abcde-fghi~jklmnopqrstuvwxyz1234567890abcdefgh~jklmnopqrstuvwxyz1234567890abcdefgh-abcd~abcdefghijklmnopqrstuvwxyz1234-567890abcde~jklmnopqrstuvwxyz1234567890abcdefgh-jklmnopqrstuvwxyz1234567890abcdefgh-567890abcde__\u0026Key-Pair-Id=ABCDEFGHIJKLMNOPQRST"}
```

5. If the book contains audio (the previous step returns something), create a folder based on JavaScript tag e.g.`"reader:book:title:changed", "Bad Blood"` on the local drive.

6. Decode to proper URL, (replace `\u0026` with `&`), e.g.:
```https://abcdefgh12345.cloudfront.net/5c28f2fc6cee070008e7a3d7/5c28f3296cee070007b46369.m4a?Expires=1234567890&Signature=abcdefghijklmnopqrstuvwxyz1234-567890abcde-fghi~jklmnopqrstuvwxyz1234567890abcdefgh~jklmnopqrstuvwxyz1234567890abcdefgh-abcd~abcdefghijklmnopqrstuvwxyz1234-567890abcde~jklmnopqrstuvwxyz1234567890abcdefgh-jklmnopqrstuvwxyz1234567890abcdefgh-567890abcde__&Key-Pair-Id=ABCDEFGHIJKLMNOPQRST```

7. Download the chapter using the above link as the m4a file. Filename will be based on `data-chapterNo` and stored in the book title folder, e.g.:
`Bad Blood/000.m4a`,
`Bad Blood/001.m4a`,
`Bad Blood/002.m4a`,... .

---
_This is currently work in progress (Thu Jan 10 14:42:44 AEDT 2019)._
