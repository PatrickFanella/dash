package domain

import "time"

type Section struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Icon        string    `json:"icon"`
	Cols        int32     `json:"cols"`
	Collapsed   bool      `json:"collapsed"`
	SortOrder   int32     `json:"sort_order"`
	SectionType string    `json:"section_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Service struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	URL            string    `json:"url"`
	Description    string    `json:"description"`
	Icon           string    `json:"icon"`
	StatusCheck    bool      `json:"status_check"`
	StatusCheckURL *string   `json:"status_check_url"`
	SortOrder      int32     `json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type NestedSection struct {
	Section
	Services []Service `json:"services"`
}
