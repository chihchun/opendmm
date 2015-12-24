package opendmm

type MovieMeta struct {
  Actresses      []string
  ActressTypes   []string
  Categories     []string
  Code           string
  CoverImage     string
  Description    string
  Directors      []string
  Genres         []string
  Label          string
  Maker          string
  MovieLength    string
  Page           string
  ReleaseDate    string
  SampleImages   []string
  Series         string
  Tags           []string
  ThumbnailImage string
  Title          string
}

func Search(query string) chan MovieMeta {
  meta := make(chan MovieMeta)
  go dmmSearch(query, meta)
  return meta
}
