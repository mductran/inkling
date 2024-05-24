package manga

import "time"

type Manga struct {
	Url          string     `bson:"url,omitempty"`
	Title        string     `bson:"title,omitempty"`
	Artist       string     `bson:"artist,omitempty"`
	Author       string     `bson:"author,omitempty"`
	Description  string     `bson:"description,omitempty"`
	Genre        string     `bson:"genre,omitempty"`
	Status       int        `bson:"status,omitempty"`
	ThumbnailUrl string     `bson:"thumbnail_url,omitempty"`
	LastUpdate   time.Time  `bson:"last_update,omitempty"`
	Chapters     *[]Chapter `bson:"chapters,omitempty"`
}

type Chapter struct {
	MangaId   string  `bson:"manga_id,omitempty"`
	Url       string  `bson:"url,omitempty"`
	Name      string  `bson:"name,omitempty"`
	Number    float32 `bson:"number,omitempty"` // float for .5 Chapter
	Scanlator string  `bson:"scanlator,omitempty"`
	Pages     []Page  `bson:"pages,omitempty"`
}

type Page struct {
	Index    int    `bson:"index,omitempty"`
	Url      string `bson:"url,omitempty"`
	ImageUrl string `bson:"image_url,omitempty"`
}
