package data

type Category struct {
	Name string
	Desc string
}

type Author struct {
	Id        string
	LastName  string
	FirstName string
}

type Book struct {
	Id         int64
	Name       string
	Authors    []*Author
	Categories []*Category
	Price      float32
	Isbn       string
}

func (book *Book) AddAuthors(authors ...*Author) {
	book.Authors = append(book.Authors, authors...)
}

func (book *Book) AddCategories(categories ...*Category) {
	book.Categories = append(book.Categories, categories...)
}

func NewBook(id int64, name string, price float32, isbn string) *Book {
	return &Book{
		Id:    id,
		Name:  name,
		Price: price,
		Isbn:  isbn,
	}
}

type Shelf struct {
	Id    string
	Name  string
	Books []*Book
}

func (shelf *Shelf) AddBooks(books ...*Book) {
	shelf.Books = append(shelf.Books, books...)
}

func NewShelf(id, name string) *Shelf {
	shelf := &Shelf{
		Id:    id,
		Name:  name,
		Books: []*Book{},
	}
	return shelf
}

func NewBookShelf() *Shelf {
	shelf := NewShelf("A", "Shelf A")

	categoryScience := &Category{
		Name: "science",
		Desc: "Science & Nature",
	}

	categoryScifi := &Category{
		Name: "scifi",
		Desc: "Science Fiction & Fantasy",
	}

	categoryHistory := &Category{
		Name: "history",
		Desc: "History",
	}

	categoryFiction := &Category{
		Name: "fiction",
		Desc: "Fiction & Literature",
	}

	bryson := &Author{
		Id:        "AF8793PL0233",
		LastName:  "Bill",
		FirstName: "Bryson",
	}

	doidge := &Author{
		Id:        "TN9984FI2323",
		LastName:  "Norman",
		FirstName: "Doidge",
	}

	shelley := &Author{
		Id:        "ER1923OW8578",
		LastName:  "Mary",
		FirstName: "Shelley",
	}

	grossman := &Author{
		Id:        "YT2223DE6644",
		LastName:  "Lev",
		FirstName: "Grossman",
	}

	susan := &Author{
		Id:        "VB9933GH3332",
		LastName:  "Bauer",
		FirstName: "Susan",
	}

	frankl := &Author{
		Id:        "LL4232FD4343",
		LastName:  "Viktor",
		FirstName: "Frankl",
	}

	winslade := &Author{
		Id:        "OW2321FF4390",
		LastName:  "William",
		FirstName: "Winslade",
	}

	kushner := &Author{
		Id:        "PW33320004",
		LastName:  "Harold",
		FirstName: "Kushner",
	}

	king := &Author{
		Id:        "AQ2324IU4343",
		LastName:  "Stephen",
		FirstName: "King",
	}

	chadbourne := &Author{
		Id:        "KU2320ZN6668",
		LastName:  "Glenn",
		FirstName: "Chadbourne",
	}

	stone := &Author{
		Id:        "MK4347NX5590",
		LastName:  "Tamara",
		FirstName: "Stone",
	}

	book := NewBook(10001, "A Short History of Nearly Everything", 12.99, "9780767908184")
	book.AddAuthors(bryson)
	book.AddCategories(categoryScience)
	shelf.AddBooks(book)

	book = NewBook(15242,
		"The Brain That Changes Itself: Stories of Personal Triumph from the Frontiers of Brain Science",
		14.98, "9780143113102")
	book.AddAuthors(doidge)
	book.AddCategories(categoryScience)
	shelf.AddBooks(book)

	book = NewBook(32998, "Frankenstein", 10.80, "9781435136168")
	book.AddAuthors(shelley)
	book.AddCategories(categoryScifi)
	shelf.AddBooks(book)

	book = NewBook(27654, "The Magician King", 5.38, "9780594465805")
	book.AddAuthors(grossman)
	book.AddCategories(categoryScifi)
	shelf.AddBooks(book)

	book = NewBook(77490,
		"The History of the Renaissance World: From the Rediscovery of Aristotle to the Conquest of Constantinople",
		21.42, "9780393059762")
	book.AddAuthors(susan)
	book.AddCategories(categoryHistory)
	shelf.AddBooks(book)

	book = NewBook(24555, "Man's Search for Meaning", 8.99, "9780807014295")
	book.AddAuthors(frankl, winslade, kushner)
	book.AddCategories(categoryHistory)
	shelf.AddBooks(book)

	book = NewBook(33989, "The Dark Man", 35.96, "9781587674259")
	book.AddAuthors(king, chadbourne)
	book.AddCategories(categoryFiction, categoryScifi)
	shelf.AddBooks(book)

	book = NewBook(18845, "Time Between Us", 9.99, "9781423159773")
	book.AddAuthors(stone)
	book.AddCategories(categoryFiction)
	shelf.AddBooks(book)

	return shelf
}
