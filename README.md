# scaneo
[![Build Status](https://drone.io/github.com/variadico/scaneo/status.png)](https://drone.io/github.com/variadico/scaneo/latest)
[![Coverage Status](https://coveralls.io/repos/variadico/scaneo/badge.svg?branch=master)](https://coveralls.io/r/variadico/scaneo?branch=master)

Generate Go code to convert database rows into arbitrary structs.
Works with any database driver. Don't have to worry about database columns
and struct names matching or tagging structs. No reflection. No ORM magic.

## Installation
If you have Go installed, then do this.
```
go get github.com/variadico/scaneo
```

Otherwise, download the standalone binary from the
[releases page](https://github.com/variadico/scaneo/releases/latest).

## Usage
```
scaneo [options] paths...
```

### Options
```
-o, -output
    Set the name of the generated file. Default is scans.go.

-p, -package
    Set the package name for the generated file. Default is current
    directory name.

-u, -unexport
    Generate unexported functions. Default is export all.

-v, -version
    Print version and exit.

-h, -help
    Print help and exit.
```

## Quick Tutorial
We start with a file that looks like this, called `tables.go`.
```go
package models

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int
	Created   time.Time
	Published pq.NullTime
	Draft     bool
	Title     string
	Body      string
}
```

Next, run `scaneo tables.go`. This will create a new file called `scans.go`.
It looks like this.
```go
// DON'T EDIT *** generated by scaneo *** DON'T EDIT //

package models

import "database/sql"

func ScanPost(r *sql.Row) (Post, error) {
	var s Post
	if err := r.Scan(
		&s.ID,
		&s.Created,
		&s.Published,
		&s.Draft,
		&s.Title,
		&s.Body,
	); err != nil {
		return Post{}, err
	}
	return s, nil
}
func ScanPosts(rs *sql.Rows) ([]Post, error) {
	structs := make([]Post, 0, 16)
	var err error
	for rs.Next() {
		var s Post
		if err = rs.Scan(
			&s.ID,
			&s.Created,
			&s.Published,
			&s.Draft,
			&s.Title,
			&s.Body,
		); err != nil {
			return nil, err
		}
		structs = append(structs, s)
	}
	return structs, nil
}
```

Now you can call those functions from other parts of your code.
```go
func serveHome(resp http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("select * from post")
	if err != nil {
		log.Println(err)
		return
	}

	posts, err := models.ScanPosts(rows) // ScanPosts was auto-generated!
	if err != nil {
		log.Println(err)
	}

	// ... send posts to template or whatever...
}
```

### Go Generate
If you want to use `scaneo` with `go generate`, then just add this comment to
the top of `tables.go`.
```go
//go:generate scaneo $GOFILE

package models
// ... rest of code...
```

Now you can call `go generate` in `package models` and `scans.go` will be
created.

## FAQ
**Why did you write this instead of using sqlx, modl, gorm, gorp, etc?**

1. Didn't know which one I should learn. Already knew `database/sql`.
2. All I wanted was structs from the database.
3. You still write SQL and declare structs with sqlx, so I didn't get the
point. Others felt unnecessary and way more complex for my small project.
4. I can SQL.

**Do my table columns have to match my struct field names?**

Nope. The names don't actually matter at all. However, what *does* matter is
the order in which you declare the types in your struct. *That* has to match
the database table. So if your table looks like this:

```
 user_id | first_name | last_name
---------+------------+-----------
       1 | Silvio     | Rodríguez
```

then, your struct field names **must** follow the same order, like this.

```go
type User struct {
	ID        int
	FirstName string
	LastName  string
}
```
