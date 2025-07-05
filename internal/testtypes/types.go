package testtypes

type User struct {
	Name string
	ID   int
}

type Post struct {
	Title  string
	Body   string
	ID     int
	UserID int
}

type UserPost struct {
	UserName string
	Body     string
	Title    string
}
