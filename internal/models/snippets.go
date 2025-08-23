package models

import (
	"database/sql"
	"time"
)

type Snippet struct{
	Id int
	Title string
	Content string
	Created time.Time
	Expired time.Time
}

type SnippetModel struct {
	DB *sql.DB
} 

func (m *SnippetModel) Insert(title string,content string, expires int)(int,error){
	return 0,nil
}

func (m *SnippetModel) Get(id int)(Snippet,error){
	return Snippet{},nil
}

func (m *SnippetModel) Latest()([]Snippet,error){
	return nil, nil
}
