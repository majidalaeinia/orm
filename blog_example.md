### Getting Started
Let's imagine we are going to build a simple blogging application that has 3 entities, `Comment`, `Post`, `Category`, `HeaderPicture`. To
start using ORM you need to call **Initialize** method. It gets array of of **ConnectionConfig** objects which has:

- `Name`: Name of the connection, it can be anything you want.
- `Driver`: Name of the driver to be used when opening connection to your database.
- `ConnectionString`: connection string to connect to your db.
- `Entities`: List of entities you want to use for that connection (later we discuss more about entities.)

```go
orm.Initialize(orm.ConnectionConfig{
    Name:             "sqlite3", // Any name
    Driver:           "sqlite3", // can be "postgres" "mysql", or any normal sql driver name
    ConnectionString: ":memory:", // Any connection string that is valid for your driver.
    Entities:         []orm.Entity{&Comment{}, &Post{}, &Category{}, &HeaderPicture{}}, // List of entities you want to use.
})
```
#### Creating database entities
Before we go further we need to talk about **Entities**, `Entity` is an interface that you ***need*** to implement for
your models/entities to let ORM work with them. So let's define our entities.

```go
package main

import "github.com/golobby/orm"

type HeaderPicture struct {
	ID     int64
	PostID int64
	Link   string
}

type Post struct {
	ID   int64
	Body string
}

type Comment struct {
	ID     int64
	PostID int64
	Body   string
}

type Category struct {
	ID    int64
	Title string
}

```
Now we need to implement `Entity` interface for each one of our database entities, so let's do it.
````go

type HeaderPicture struct {
	ID     int64
	PostID int64
	Link   string
}

func (h HeaderPicture) ConfigureEntity(e *orm.EntityConfigurator) {
	e.Table("header_pictures")
}

func (h HeaderPicture) ConfigureRelations(r *orm.RelationConfigurator) {
	r.BelongsTo(Post{}, orm.BelongsToConfig{})
}

type Post struct {
	ID   int64
	Body string
}

func (p Post) ConfigureEntity(e *orm.EntityConfigurator) {
	e.
		Table("posts")

}

func (p Post) ConfigureRelations(r *orm.RelationConfigurator) {
	r.
		HasMany(Comment{}, orm.HasManyConfig{}).
		HasOne(HeaderPicture{}, orm.HasOneConfig{}).
		BelongsToMany(Category{}, orm.BelongsToManyConfig{IntermediateTable: "post_categories"})
}

func (p *Post) Categories() ([]Category, error) {
	return orm.BelongsToMany[Category](p)
}

func (p *Post) Comments() ([]Comment, error) {
	return orm.HasMany[Comment](p)
}

type Comment struct {
	ID     int64
	PostID int64
	Body   string
}

func (c Comment) ConfigureEntity(e *orm.EntityConfigurator) {
	e.Table("comments")
}

func (c Comment) ConfigureRelations(r *orm.RelationConfigurator) {
	r.BelongsTo(Post{}, orm.BelongsToConfig{})
}

func (c *Comment) Post() (Post, error) {
	return orm.BelongsTo[Post](c)
}

type Category struct {
	ID    int64
	Title string
}

func (c Category) ConfigureEntity(e *orm.EntityConfigurator) {
	e.Table("categories")
}

func (c Category) ConfigureRelations(r *orm.RelationConfigurator) {
	r.BelongsToMany(Post{}, orm.BelongsToManyConfig{IntermediateTable: "post_categories"})
}

func (c Category) Posts() ([]Post, error) {
	return orm.BelongsToMany[Post](c)
}

````
As you can see each `Entity` should define 2 methods: `ConfigureEntity` which basically maps the struct to correct table and database connection and `ConfigureRelations` registers it's relations with other models.

after creating all entities you can validate your entities using `orm.Schematic()` method to see all information that ORM extracted and constructed from your code.
```text
----------------default---------------
SQL Dialect: sqlite3
Table: categories
+----------+--------+----------------+------------+
| SQL NAME | TYPE   | IS PRIMARY KEY | IS VIRTUAL |
+----------+--------+----------------+------------+
| id       | int64  | true           | false      |
| title    | string | false          | false      |
+----------+--------+----------------+------------+
categories N-N posts => {IntermediateTable:post_categories IntermediatePropertyID:post_id IntermediateOwnerID:category_id ForeignTable:posts ForeignLookupColumn:id}

Table: header_pictures
+----------+--------+----------------+------------+
| SQL NAME | TYPE   | IS PRIMARY KEY | IS VIRTUAL |
+----------+--------+----------------+------------+
| id       | int64  | true           | false      |
| post_id  | int64  | false          | false      |
| link     | string | false          | false      |
+----------+--------+----------------+------------+
header_pictures N-1 posts => {OwnerTable:posts LocalForeignKey:post_id ForeignColumnName:id}

Table: comments
+----------+--------+----------------+------------+
| SQL NAME | TYPE   | IS PRIMARY KEY | IS VIRTUAL |
+----------+--------+----------------+------------+
| id       | int64  | true           | false      |
| post_id  | int64  | false          | false      |
| body     | string | false          | false      |
+----------+--------+----------------+------------+
comments N-1 posts => {OwnerTable:posts LocalForeignKey:post_id ForeignColumnName:id}

Table: posts
+----------+--------+----------------+------------+
| SQL NAME | TYPE   | IS PRIMARY KEY | IS VIRTUAL |
+----------+--------+----------------+------------+
| id       | int64  | true           | false      |
| body     | string | false          | false      |
+----------+--------+----------------+------------+
posts 1-N comments => {PropertyTable:comments PropertyForeignKey:post_id}
posts 1-1 header_pictures => {PropertyTable:header_pictures PropertyForeignKey:post_id}
posts N-N categories => {IntermediateTable:post_categories IntermediatePropertyID:category_id IntermediateOwnerID:post_id ForeignTable:categories ForeignLookupColumn:id}
```

#### Create, Find, Update, Delete
Now let's write simple `CRUD` logic for posts.

```go
package main

import "github.com/golobby/orm"

func createPost(p *Post) error {
	err := orm.Save(p)
	return err
}
func findPost(id int) (*Post, error) {
	return orm.Find[Post](id)
}

func updatePost(p *Post) error {
	return orm.Update(p)
}

func deletePost(p *Post) error {
	return orm.Delete(p)
}

```
#### Insert with relation
now that we have our post in database, let's add some comments to it. notice that comments are in relation with posts and the relation from posts view is a hasMany relationship and from comments is a belongsTo relationship.

```go
package main

func addCommentsToPost(post *Post, comments []Comment) error {
	return orm.Add(post, comments)
}

func addComments(comments []Comment) error {
	return orm.SaveAll(comments)
}

// you can also create, update, delete, find comments like you saw with posts.
```

finally, now we have both our posts and comments in db, let's add some categories.

```go
package main

func addCategoryToPost(post *Post, category *Category) error {
	return orm.Add(post, category)
}


```
#### Custom query
Now what if you want to do some complex query for example to get some posts that were created today ?

```go
package main

import "github.com/golobby/orm"

func getTodayPosts() ([]Post, error) {
	posts, err := orm.Query[Post](
		orm.
			Select().
			Where("created_at", "<", "NOW()").
			Where("created_at", ">", "TODAY()").
			OrderBy("id", "desc"))
    return posts, err
}
```
basically you can use all orm power to run any custom query, you can build any custom query using orm query builder but you can even run raw queries and use orm power to bind them to your entities.
You can see querybuilder docs in [query builder package](https://github.com/golobby/orm/tree/master/querybuilder)
```go
package main

import "github.com/golobby/orm"

func getTodayPosts() ([]Post, error) {
	return orm.RawQuery[Post]("SELECT * FROM posts WHERE created_at < NOW() and created_at > TODAY()")
}
```