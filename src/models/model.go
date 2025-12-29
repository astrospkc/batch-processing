package models

import "time"

type User struct{
	Id		string 		`bson:"id,omitempty" json:"id"`
	Name	string		`bson:"name" json:"name"`
	Email   string		`bson:"email" json:"email"` // should be unique
	Posts   []Post		`bson:"posts,omitempty" json:"posts"`
	API_Key string		`bson:"api_key" json:"api_key"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`   
}

type Post struct{
	Id 		string 		`bson:"id,omitempty" json:"id"`
	UserId  string		`bson:"user_id" json:"user_id"`
	Title	string		`bson:"title" json:"title"`
	Content	string 		`bson:"content" json:"content"`
	LikeCount int64     `bson:"like_count" json:"like_count"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Likes struct{
	Id		string 		`bson:"id,omitempty" json:"id"`
	PostId	string		`bson:"post_id" json:"post_id"`
	UserId  string      `bson:"user_id" json:"user_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

}
