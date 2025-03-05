package testtypes

type User struct {
	ID   int
	Name string
}

type Post struct {
	ID     int
	UserID int
	Title  string
	Body   string
}

type UserPost struct {
	UserName string
	Body     string
	Title    string
}
