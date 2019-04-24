package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	api "sample-grpc/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	db *sql.DB
}

// NewBookServiceServer return new server
func NewBookServiceServer(db *sql.DB) api.BookServiceServer {
	return &server{db: db}
}

func (s *server) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	book := req.GetBook()
	res := api.CreateResponse{}

	err := s.db.QueryRow(
		"insert into books (title, author, description, pages, price) values ($1, $2, $3, $4, $5) returning id",
		book.GetTitle(), book.GetAuthor(), book.GetDescription(), book.GetPages(), book.GetPrice(),
	).Scan(&res.Id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &res, nil
}
func (s *server) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	id := req.GetId()
	rows, err := s.db.Exec("delete from books where id=$1", id)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	countRows, err := rows.RowsAffected()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if countRows != 1 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID='%d' is not found", req.Id))
	}

	return &api.DeleteResponse{
		Deleted: countRows,
	}, nil
}

func (s *server) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponce, error) {
	book := req.GetBook()
	rows, err := s.db.Exec(
		"update books set title=$1, author=$2, description=$3, pages=$4, price=$5 where id=$6",
		book.GetTitle(), book.GetAuthor(), book.GetDescription(), book.GetPages(), book.GetPrice(), book.GetId())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	countRows, err := rows.RowsAffected()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if countRows != 1 {
		return nil, err
	}

	return &api.UpdateResponce{
		Updated: countRows,
	}, nil
}

func (s *server) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	id := req.GetId()
	row, err := s.db.Query("select id, title, author, description, pages, price from books where id=$1", id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer row.Close()

	book := api.Book{}
	err = row.Scan(&book.Id, &book.Title, &book.Author, &book.Description, &book.Pages, &book.Price)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	res := api.GetResponse{}
	res.Book = &book
	return &res, nil
}

func (s *server) GetAll(ctx context.Context, req *api.GetAllRequest) (*api.GetAllResponse, error) {
	rows, err := s.db.Query("select id, title, author, description, pages, price from books")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	list := []*api.Book{}
	for rows.Next() {
		book := api.Book{}

		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Description, &book.Pages, &book.Price)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		list = append(list, &book)
	}

	res := api.GetAllResponse{}
	res.Books = list
	return &res, nil
}
