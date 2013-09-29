package data

type User struct {
	Id        int64
	LoginName string
	DispName  string
}

type Comment struct {
	Id      int64
	Owner   *User
	Content string
}

type Post struct {
	Id       int64
	Owner    *User
	BlogId   int64
	Content  string
	Comments []*Comment
}

type Blog struct {
	Id    int64
	Owner *User
	Posts []*Post
}

func NewUser(id int64, loginName string, dispName string) *User {
	return &User{
		Id:        id,
		LoginName: loginName,
		DispName:  dispName,
	}
}

func NewBlog() *Blog {
	owner := NewUser(1, "johndoe", "John Doe")
	visitor1 := NewUser(2, "joeblow", "Joe Blow")
	visitor2 := NewUser(3, "janeroe", "Jane Roe")

	return &Blog{
		Id:    10010001,
		Owner: owner,
		Posts: []*Post{
			&Post{
				Id:      50001001,
				Owner:   owner,
				BlogId:  10010001,
				Content: "I like eating donuts.",
				Comments: []*Comment{
					&Comment{
						Id:      90001001,
						Owner:   visitor1,
						Content: "Me too!",
					},
					&Comment{
						Id:      90001002,
						Owner:   visitor2,
						Content: "Not healthy.",
					},
					&Comment{
						Id:      90001003,
						Owner:   visitor1,
						Content: "Come on!",
					},
				},
			},
			&Post{
				Id:      50001002,
				Owner:   owner,
				BlogId:  10010001,
				Content: "Jee, this morning I got up late.",
				Comments: []*Comment{
					&Comment{
						Id:      90002001,
						Owner:   visitor2,
						Content: "Ah oh! You are saying you missed the exam?",
					},
					&Comment{
						Id:      90002002,
						Owner:   visitor1,
						Content: "Crying for you, darling. :)",
					},
					&Comment{
						Id:      90002003,
						Owner:   visitor2,
						Content: "I hope he is joking!",
					},
					&Comment{
						Id:      90002002,
						Owner:   visitor1,
						Content: "Well, I think it's a punishment for always get to bed at 3am.",
					},
				},
			},
		},
	}
}
