package models

type User struct {
	ID                int
	TelegramID        int64
	Username          string
	Role              string
	Timezone          string
	NotificationTime  int
	NotificationCount int
}

type Source struct {
	ID       int
	Name     string
	URL      string
	Category string
}

type News struct {
	ID        int
	SourceID  int
	Tilte     string
	Link      string
	Published string
	Category  string
}
