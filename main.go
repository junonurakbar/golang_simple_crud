package main

// tambahin fitur update, dan delete sendiri!

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Book struct {
	Id, Pages                  int
	Title, Author, ReleaseYear string
}

var books []Book
var fileName = "data.csv"  // file data formatnya CSV

func main() {
	createFile(fileName)
	LoadDataFromCSV(fileName)

	for {
		fmt.Println(strings.Repeat("=", 50))
		fmt.Print(strings.Repeat("=", 15))
		fmt.Print("Books Data Management")
		fmt.Println(strings.Repeat("=", 14))
	
		fmt.Println("1. View All Books")
		fmt.Println("2. Add a New Book")
		fmt.Println("3. Update a Book")
		fmt.Println("4. Delete a Book")
		fmt.Println("5. Exit")
		fmt.Println(strings.Repeat("=", 50))
		
		fmt.Print("Enter your choice: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := scanner.Text()
		if strings.Contains(answer, "1") || strings.Contains(answer, strings.ToLower("View")) {
			err := viewAllBooks()
			if err != nil { fmt.Println(err) }
			
		} else if strings.Contains(answer, "2") || strings.Contains(answer, strings.ToLower("Add")) {
			err := addNewBook()
			if err != nil { fmt.Println(err) }

		} else if strings.Contains(answer, "3") || strings.Contains(answer, strings.ToLower("Update")) {
			err := UpdateBook()
			if err != nil { fmt.Println(err) }

		} else if strings.Contains(answer, "4") || strings.Contains(answer, strings.ToLower("Delete")) {
			err := DeleteBook()
			if err != nil { fmt.Println(err) }

		} else {
			fmt.Println("Good bye!")
			os.Exit(0)
		}
	}

	// addNewBook()
	// UpdateBook()
	// DeleteBook()
	// LoadDataFromCSV(fileName)
	// viewAllBooks()
}

func createFile(fileName string) {
	// create file
	_, err := os.Stat(fileName) // returns file info describing the named file
	if os.IsNotExist(err) { // jika file tidak exist
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Println("File", fileName, "berhasil dibuat")
	}
}

func addNewBook() error {
	var newBook Book
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter Book Details")

	fmt.Print("Book Id: ")
	scanner.Scan()
	newBook.Id, _ = strconv.Atoi(scanner.Text())
	_, err := FindBookById(newBook.Id)
	if err == nil {
		return fmt.Errorf("book with id: %d already exists", newBook.Id)
	}

	fmt.Print("Book Title: ")
	scanner.Scan()
	newBook.Title = scanner.Text()

	fmt.Print("Book Author: ")
	scanner.Scan()
	newBook.Author = scanner.Text()

	fmt.Print("Book ReleaseYear: ")
	scanner.Scan()
	newBook.ReleaseYear = scanner.Text()

	fmt.Print("Book Pages: ")
	scanner.Scan()
	newBook.Pages, _ = strconv.Atoi(scanner.Text())

	_, err = FindBookById(newBook.Id)
	if err != nil {
		books = append(books, newBook)
		} else {
		return fmt.Errorf("book with id: %d already exists", newBook.Id)
	}

	fmt.Print("Are you sure you want to add this book? (y/n) ")
	scanner.Scan()
	answer := scanner.Text()
	if strings.Contains(answer, strings.ToLower("y")) {
		err = SaveDataToCSV(fileName)
		if err != nil {
			return err
		}
		fmt.Println("book added successfully")
	} else {
		fmt.Println("book canceled to be added")
	}
	return nil
}

func SaveDataToCSV(fileName string) error {
	// file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0666)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error opening csv file: %w", err)
	}
	defer file.Close()

	for _, book := range books{
		// row: data yang akan ditulis di file csv
		row := strconv.Itoa(book.Id) + "," + book.Title + "," + book.Author + "," + book.ReleaseYear + "," + strconv.Itoa(book.Pages) + "\n"
		file.WriteString(row)
	}
	return nil
}

func FindBookById(id int) (Book, error) {
	for _, book := range books {
		if book.Id == id {
			return book, nil
		}
	}
	return Book{}, fmt.Errorf("id: %d not found", id)
}

func LoadDataFromCSV(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error open csv file %w", err)
	}
	defer file.Close()

	// baca isi data, read line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// record: utk menampung file data dari csv
		record := strings.Split(scanner.Text(), ",") // read one line from csv and then split them by the comma, after that insert them into a slice
		// fmt.Println(record)

		// masukkin semua data yang ada di variabel record ke variabel baru book yg jenisnya struct "Book"
		id, _ := strconv.Atoi(record[0])
		pages, _ := strconv.Atoi(record[4])
		book := Book{
			Id: id,
			Title: record[1],
			Author: record[2],
			ReleaseYear: record[3],
			Pages: pages,
		}
		books = append(books, book)
	}
	// cek apabila ada error saat scanning file csv
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("error opening csv file %w", err)
	}
	return nil
}

func UpdateBook() error {
	// masukin id buku yang mau diganti detailnya
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Print("Enter book id that you want to change: ")
	scanner.Scan()
	findId, _ := strconv.Atoi(scanner.Text())
	currentBook, _ := FindBookById(findId)
	newBook, err := FindBookById(findId)
	if err != nil {
		return fmt.Errorf("error, book not found %w", err)
	} else {
		fmt.Print("Enter book title that you want to change: ")
		scanner.Scan()
		newBook.Title = scanner.Text()
		if strings.Trim(newBook.Title, " ") == "" {
			newBook.Title = currentBook.Title
		}

		fmt.Print("Enter book author that you want to change: ")
		scanner.Scan()
		newBook.Author = scanner.Text()
		if strings.Trim(newBook.Author, " ") == "" {
			newBook.Author = currentBook.Author
		}

		fmt.Print("Enter book releaseYear that you want to change: ")
		scanner.Scan()
		newBook.ReleaseYear = scanner.Text()
		if strings.Trim(newBook.ReleaseYear, " ") == "" {
			newBook.ReleaseYear = currentBook.ReleaseYear
		}

		fmt.Print("Enter book pages that you want to change: ")
		scanner.Scan()
		newBook.Pages, _ = strconv.Atoi(scanner.Text())
		if newBook.Pages == 0 || scanner.Text() == "" {
			newBook.Pages = currentBook.Pages
		}
		
		fmt.Printf("Book Id: %v\n", newBook.Id)
		fmt.Printf("Book Title: %v\n", newBook.Title)
		fmt.Printf("Book Author: %v\n", newBook.Author)
		fmt.Printf("Book ReleaseYear: %v\n", newBook.ReleaseYear)
		fmt.Printf("Book Pages: %d\n", newBook.Pages)
		// confirm if you want to update the book information
		fmt.Print("Are you sure you want to update the book? (y/n) ")
		scanner.Scan()
		answer := scanner.Text()
		if strings.Contains(answer, strings.ToLower("y")) {
			// update books from the slice, but its still not updated to the csv file
			for i, book := range books {	
				if book.Id == newBook.Id {
					books[i] = newBook
				}
			}

			// overwrite the csv file
			file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR, 0666)
			if err != nil {
				return fmt.Errorf("unable to update a book, %w", err)
			}
			for _, book := range books{
				// row: data yang akan ditulis di file csv
				row := strconv.Itoa(book.Id) + "," + book.Title + "," + book.Author + "," + book.ReleaseYear + "," + strconv.Itoa(book.Pages) + "\n"
				file.WriteString(row)
			}
			fmt.Println("Book updated successfully")
		} else {
			fmt.Println("Book canceled to be updated")
		}
	}

	return nil
}

func DeleteBook() error {
	// masukin id buku yang mau diganti detailnya
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Print("Enter book id that you want to delete: ")
	scanner.Scan()
	findId, _ := strconv.Atoi(scanner.Text())
	foundBook, err := FindBookById(findId)
	if err != nil {
		return fmt.Errorf("error, book not found %w", err)
	} else {
		showBookById(findId)
		fmt.Print("Are you sure you want to delete this book? (yes/no) ")
		scanner.Scan()
		answer := scanner.Text()
		// delete one book from the slice
		if strings.Contains(answer, strings.ToLower("y")) {
			// cari index dari slice "books", buat ngehapus element yang diingingkan
			// 		 slices.Index(slice, element_dari_slice_yang_dicari)
			index := slices.Index(books, foundBook)	
			// hapus element yang diinginkan berdasarkan variabel "index" dari slice "books"
			books = append(books[:index], books[index+1:]...)
			
			// overwrite the csv file
			file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR, 0666)
			if err != nil {
				return fmt.Errorf("unable to update a book, %w", err)
			}
			for _, book := range books{
				// row: data yang akan ditulis di file csv
				row := strconv.Itoa(book.Id) + "," + book.Title + "," + book.Author + "," + book.ReleaseYear + "," + strconv.Itoa(book.Pages) + "\n"
				file.WriteString(row)
			}
			fmt.Println("book deleted successfully")
		} else {
			fmt.Println("book canceled to be deleted")
		}
	}
	return nil
}

func viewAllBooks() error {
	if len(books) == 0 {	//jika slice books tidak mempunyai elemen
		return fmt.Errorf("no books available")
	}

	for i, book := range books {
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("Book - %d\n", i+1)
		fmt.Printf("Book Id: %v\n", book.Id)
		fmt.Printf("Book Title: %v\n", book.Title)
		fmt.Printf("Book Author: %v\n", book.Author)
		fmt.Printf("Book ReleaseYear: %v\n", book.ReleaseYear)
		fmt.Printf("Book Pages: %d\n", book.Pages)
		fmt.Println(strings.Repeat("=", 50))
	}
	return nil
}

func showBookById(id int) error {
	if len(books) == 0 {
		return fmt.Errorf("no books available")
	}

	searchBook, err := FindBookById(id)
	if err != nil {
		return fmt.Errorf("error, no book with id %d", id)
	}

	fmt.Printf("Book Id: %v\n", searchBook.Id)
	fmt.Printf("Book Title: %v\n", searchBook.Title)
	fmt.Printf("Book Author: %v\n", searchBook.Author)
	fmt.Printf("Book ReleaseYear: %v\n", searchBook.ReleaseYear)
	fmt.Printf("Book Pages: %d\n", searchBook.Pages)

	return nil
}