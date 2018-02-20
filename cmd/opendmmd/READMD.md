# OpenDMM Server

## Usage

### Search API

```bash
$ curl -i "http://api.libredmm.com/search?q=ABP-123"
HTTP/1.1 200 OK
Server: Cowboy
Connection: keep-alive
Date: Tue, 06 Feb 2018 07:54:42 GMT
Content-Length: 649
Content-Type: text/plain; charset=utf-8
Via: 1.1 vegur

{
  "Actresses": [
    "酒井ももか"
  ],
  "ActressTypes": null,
  "Categories": null,
  "Code": "ABP-123",
  "CoverImage": "http://pics.dmm.co.jp/mono/movie/adult/118abp123/118abp123pl.jpg",
  "Description": "",
  "Directors": [
    "Porn Stars"
  ],
  "Genres": [
    "単体作品",
    "巨乳",
    "風俗嬢"
  ],
  "Label": "ABSOLUTELY PERFECT",
  "Maker": "プレステージ",
  "MovieLength": "140",
  "Page": "http://www.javlibrary.com/ja/?v=javlijqo4e",
  "ReleaseDate": "2014-04-01",
  "SampleImages": null,
  "Series": "",
  "Tags": null,
  "ThumbnailImage": "",
  "Title": "酒井ももか、満足度満点新人ソープ DX"
}
```

```
$ curl -i "http://api.libredmm.com/search?q=NONSENSE-789"
HTTP/1.1 404 Not Found
Server: Cowboy
Connection: keep-alive
Date: Tue, 06 Feb 2018 07:56:05 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Via: 1.1 vegur
```

### Guess API

```
$ curl -i "http://api.libredmm.com/guess?q=ABP-123+test"
HTTP/1.1 200 OK
Server: Cowboy
Connection: keep-alive
Date: Tue, 06 Feb 2018 07:58:51 GMT
Content-Length: 15
Content-Type: text/plain; charset=utf-8
Via: 1.1 vegur

[
  "ABP-123"
]
```

```
$ curl -i "http://api.libredmm.com/guess?q=NONSENSE"
HTTP/1.1 404 Not Found
Server: Cowboy
Connection: keep-alive
Date: Tue, 06 Feb 2018 07:57:03 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Via: 1.1 vegur
```
