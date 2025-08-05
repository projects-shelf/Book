# shelf | Book

### A lightweight ebook server with an intuitive UI

<br>

![web client](./assets/shelf_book_web.png)

### Supported Formats

| Format     | Viewer | Extract Cover | Extract Metadata |
|------------|--------|----------------|------------------|
| PDF        | ✅      | ✅             | ✅               |
| EPUB       | ✅      | ✅             | ✅               |
| CBZ        | ✅      | ✅             | ❌               |
| CBR        | ✅      | ✅             | ❌               |

✅ = Supported  △ = Partial Support / Experimental  ❌ = Not Supported

<details>
<summary><strong>Known Bugs</strong></summary>

- Safari: The epubViewer's display area may be reduced.  
- Safari: The epubViewer may become unresponsive.  

</details>

## Quickstart

```shell
git clone https://github.com/projects-shelf/Book.git
cd Book
docker-compose up -d
```

During the initial scan of the books, the server may appear unresponsive. Please wait until the process is finished.

## Features

### Simple viewer

![viewer](./assets/shelf_book_viewer.png)

options
  - LtoR / RtoL reading direction
  - No spreads / Odd spreads / Even spreads
  - Font size

### Searchable by title and metadata

![search](./assets/shelf_book_search.png)

Allows prefix search on titles and exact match search on metadata using hashtags (#).

### Lightweight

Built with Go and SQLite, it consumes minimal system resources.

```
CONTAINER ID   NAME                 CPU %     MEM USAGE / LIMIT     MEM %     NET I/O         BLOCK I/O   PIDS
13ac41d1866e   shelf_book_nginx     0.00%     1.91MiB / 3.88GiB    0.05%     209MB / 210MB     11.5MB / 541kB    2
bcebb5552894   shelf_book_go        0.00%     14.57MiB / 3.88GiB   0.37%     706kB / 208MB     111MB / 38.7MB    11
```

Resource consumption during idle periods (on Intel Mac)

## Note

- PDF rendering requires a moderate amount of resources.

## License

Book is licensed under [MIT License](https://github.com/projects-shelf/Book/blob/main/LICENSE).

## Author

Developed by [PepperCat](https://github.com/PepperCat-YamanekoVillage).
