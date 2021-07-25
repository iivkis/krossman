package adidasParser

type Config struct {
	SaveAs     string
	Address    string
	User       string
	Code       string
	FeedID     int
	CategoryID string
}

type pr struct {
	limit    int
	offset   int
	filepath string
	config   *Config
}
