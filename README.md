# OpenDMM

[![Build Status](https://travis-ci.org/libredmm/opendmm.svg?branch=master)](https://travis-ci.org/libredmm/opendmm)
[![codecov](https://codecov.io/gh/libredmm/opendmm/branch/master/graph/badge.svg)](https://codecov.io/gh/libredmm/opendmm)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/libredmm/opendmm/blob/master/LICENSE)

This Go Librady/CLI fetches AV metadata given a code (番号).

## Installation & Usage

    $ go install github.com/libredmm/opendmm/cmd/opendmm
    $ opendmm ABP-123
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

It's also available as a library. See [the CLI code](https://github.com/libredmm/opendmm/blob/master/cmd/opendmm/opendmm.go) for example usage.
