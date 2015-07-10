// Wanna See It? API example, by Saam Aghevli

package main

import (
	"crypto/md5"
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// Data structure for users. Full OAUTH integration would be preferred in prod.
type User struct {
	Uid int64 `db:"Uid"`
	Pwd string
}

// Data structure for a full WSI post. Note that the image data is simply
// mocked out as a string.
type Post struct {
	Post_id int64 `db:"Pid"`
	Text    string
	Author  int64
	Date    int64
	Img     string
}

// Simper data structure that is returned for post permalinks.
type SimplePost struct {
	Text string
	Img  string
}

// Globals are gross, but this makes my interfaces much cleaner.
var dbmap = initDb()

func main() {
	// Could pass in a deterministic seed for testing purposes.
	rand.Seed(time.Now().UnixNano())
	defer dbmap.Db.Close()

	router := gin.Default()
	router.POST("/login", CreateUser)
	router.POST("/postimg", CreatePost)
	router.GET("/posts/:count/:offset", PostsList)
	router.GET("/perma/:id", GetPermalink)
	router.Run(":8000")
}

// POST method to create a new user. Full CRUD functionality needed in a proper
// API.
func CreateUser(gin_context *gin.Context) {
	var input User

	gin_context.Bind(&input)
	err, new_user := CreateDbUser(input.Uid, input.Pwd)
	if err == nil {
		// Password has would not normally be returned, doing so for demonstration.
		content := gin.H{
			"result":        "Success",
			"user_id":       string(new_user.Uid),
			"password_hash": new_user.Pwd,
		}
		gin_context.JSON(201, content)
	} else {
		gin_context.JSON(500, gin.H{"result": "An error occured"})
	}
}

// Handles DB interaction for user creation. Hashes given password string using
// MD5, although further crypto would be needed for a production system.
func CreateDbUser(input_user_id int64, input_pwd string) (error, User) {
	// Password should be hashed. In a production system, a random, deterministic
	// salt should be used to prevent replay attacks.
	hash := md5.New()
	io.WriteString(hash, input_pwd)
	db_user := User{
		Uid: input_user_id,
		Pwd: string(hash.Sum(nil)),
	}
	err := dbmap.Insert(&db_user)
	// Ignore return value, error is handled in caller.
	CheckAndLogError(err, "Error adding user to users table")
	return err, db_user
}

// POST method to create a new WSI post.
func CreatePost(gin_context *gin.Context) {
	var input Post
	gin_context.Bind(&input)
	// Author is currently just a stub, could be used to verify identity, for
	// example using an OAUTH token.
	err, new_post :=
		CreateDbPost(input.Author, input.Text, input.Img)
	if err == nil {
		content := gin.H{
			"result":  "Success",
			"post_id": strconv.FormatInt(new_post.Post_id, 10),
		}
		gin_context.JSON(201, content)
	} else {
		gin_context.JSON(500, gin.H{"result": "An error occured"})
	}
}

// Handles DB interaction for post creation. Assigns a random 63 bit integer
// as a simple "UUID", but more robust checks would be needed in production.
func CreateDbPost(input_author int64, input_text,
	input_img string) (error, Post) {
	// Simplified, language supported method of creating UUID's. In production,
	// a more robust solution for random post identifiers should be used. In this
	// prototype, 2^63 (9 quintillion) possibilities will suffice. Unless it goes
	// really, really viral.
	db_post := Post{
		Post_id: rand.Int63(),
		Text:    input_text,
		Author:  input_author,
		Date:    time.Now().UnixNano(),
		Img:     input_img,
	}
	err := dbmap.Insert(&db_post)
	// Ignore return value, error is handled in caller.
	CheckAndLogError(err, "Error adding post to posts table")
	return err, db_post
}

// GET method for post consuming, returns a list of the K most recent posts.
// Defaults to 5 if post count fails to parse.
// A full API could also specify different ways to consume the post - by user,
// by text, etc.
func PostsList(gin_context *gin.Context) {
	count := gin_context.Params.ByName("count")
	offset := gin_context.Params.ByName("offset")
	post_count, err_count := strconv.Atoi(count)
	post_offset, err_offset := strconv.Atoi(offset)
	if err_count != nil || err_offset != nil {
		log.Println("Error processing inputs. Defaulting to 5 posts, no offset.")
		post_count = 5
		post_offset = 0
	}
	var posts []Post
	_, err_db := dbmap.Select(
		&posts, "select * from posts order by Date limit ?, ?", post_offset,
		post_count)
	CheckAndLogError(err_db, "Select failed")
	content := gin.H{}
	for key, value := range posts {
		content[strconv.Itoa(key)] = value
	}
	gin_context.JSON(200, content)
}

// GET method for permalinks. Posts are available at the permalink specified by
// their UUID.
func GetPermalink(gin_context *gin.Context) {
	post_id := gin_context.Params.ByName("id")
	p_id, err := strconv.ParseInt(post_id, 10, 64)
	if err != nil {
		gin_context.JSON(400, gin.H{"result": "An error occured, malformed ID"})
		return
	}
	simple_post, err := GetPermalinkDb(p_id)
	if err != nil {
		gin_context.JSON(400, gin.H{"result": "Error retrieving image"})
		return
	}
	content := gin.H{"text": simple_post.Text, "image": simple_post.Img}
	gin_context.JSON(200, content)
}

// Handles DB interaction for permalinks. Returns a simpler post - just the
// text and accompanying image.
func GetPermalinkDb(p_id int64) (SimplePost, error) {
	post := Post{}
	err := dbmap.SelectOne(&post, "select * from posts where Pid=?", p_id)
	CheckAndLogError(err, "Post search failed")
	simple_post := SimplePost{
		Text: post.Text,
		Img:  post.Img,
	}
	return simple_post, err
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "wsi_db.bin")
	CheckAndFailError(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Uid")
	// Let table choose keys randomly.
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(false, "Pid")
	err = dbmap.CreateTablesIfNotExists()
	CheckAndFailError(err, "Create tables failed")
	return dbmap
}

// If there is an error, the program fails.
func CheckAndFailError(err error, message string) error {
	if err != nil {
		log.Fatalln(message, err)
	}
	return err
}

// If there is an error, logs the passed in message and returns true. Otherwise,
// returns false.
func CheckAndLogError(err error, message string) bool {
	if err != nil {
		log.Println(message, err)
		return true
	}
	return false
}
