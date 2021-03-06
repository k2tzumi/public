// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models // import "cirello.io/bookmarkd/pkg/models"

import (
	"net/url"
	"time"

	"cirello.io/errors"
	"github.com/jmoiron/sqlx"
)

// Bookmark stores the basic information of a web URL.
type Bookmark struct {
	ID               int64     `db:"id" json:"id"`
	URL              string    `db:"url" json:"url"`
	LastStatusCode   int64     `db:"last_status_code" json:"last_status_code"`
	LastStatusCheck  int64     `db:"last_status_check" json:"last_status_check"`
	LastStatusReason string    `db:"last_status_reason" json:"last_status_reason"`
	Title            string    `db:"title" json:"title"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	Inbox            int64     `db:"inbox" json:"inbox"`

	Host string `db:"-" json:"host"`
}

// BookmarkDAO provides DB persistence to bookmarks.
type BookmarkDAO interface {
	All() ([]*Bookmark, error)
	Bootstrap() error
	Delete(*Bookmark) error
	Expired() ([]*Bookmark, error)
	GetByID(id int64) (*Bookmark, error)
	Insert(*Bookmark) (*Bookmark, error)
	Invalid() ([]*Bookmark, error)
	Update(*Bookmark) error
}

type bookmarkDAO struct {
	db *sqlx.DB
}

// NewBookmarkDAO instanties a BookmarkDAO
func NewBookmarkDAO(db *sqlx.DB) BookmarkDAO {
	return &bookmarkDAO{db: db}
}

// Bootstrap creates table if missing.
func (b *bookmarkDAO) Bootstrap() error {
	cmds := []string{
		`create table if not exists bookmarks (
			id integer primary key autoincrement,
			url text,
			last_status_code int,
			last_status_check int,
			last_status_reason text,
			title bigtext not null,
			created_at datetime not null,
			inbox bool not null default false
		);
		`,
		`create index if not exists bookmarks_last_status_code  on bookmarks (last_status_code)`,
		`create index if not exists bookmarks_last_status_check on bookmarks (last_status_check)`,
		`create index if not exists bookmarks_created_at on bookmarks (created_at)`,
		`create index if not exists bookmarks_inbox on bookmarks (inbox)`,
	}

	for _, cmd := range cmds {
		_, err := b.db.Exec(cmd)
		if err != nil {
			return errors.E(err)
		}
	}

	return nil
}

// All returns all known bookmarks.
func (b *bookmarkDAO) All() ([]*Bookmark, error) {
	var bookmarks []*Bookmark
	err := b.db.Select(&bookmarks, `
		SELECT
			*

		FROM
			bookmarks

		ORDER BY
			CASE
				WHEN last_status_code = 0 THEN 999
				ELSE last_status_code
			END ASC,
			created_at DESC,
			id DESC
	`)

	for _, b := range bookmarks {
		u, err := url.Parse(b.URL)
		if err == nil {
			b.Host = u.Host
		}
	}
	return bookmarks, errors.E(err)
}

// Expired return all valid but expired bookmarks.
func (b *bookmarkDAO) Expired() ([]*Bookmark, error) {
	var bookmarks []*Bookmark
	const week = 7 * 24 * time.Hour
	deadline := time.Now().Add(-week).Unix()
	err := b.db.Select(&bookmarks, `
		SELECT
			*
		FROM
			bookmarks
		WHERE
			last_status_code = 200
			AND last_status_check <= $1
	`, deadline)
	for _, b := range bookmarks {
		u, err := url.Parse(b.URL)
		if err == nil {
			b.Host = u.Host
		}
	}
	return bookmarks, errors.E(err)
}

// Invalid return all invalid  bookmarks.
func (b *bookmarkDAO) Invalid() ([]*Bookmark, error) {
	var bookmarks []*Bookmark
	err := b.db.Select(&bookmarks, `
		SELECT
			*
		FROM
			bookmarks
		WHERE
			last_status_code != 200
	`)
	if err != nil {
		return nil, errors.E(err)
	}
	for _, b := range bookmarks {
		u, err := url.Parse(b.URL)
		if err == nil {
			b.Host = u.Host
		}
	}
	return bookmarks, nil
}

// Insert one bookmark.
func (b *bookmarkDAO) Insert(bookmark *Bookmark) (*Bookmark, error) {
	bookmark.CreatedAt = time.Now()
	bookmark.Inbox = 1
	result, err := b.db.NamedExec(`
		INSERT INTO bookmarks
		(url, last_status_code, last_status_check, last_status_reason, title, created_at, inbox)
		VALUES
		(:url, :last_status_code, :last_status_check, :last_status_reason, :title, :created_at, :inbox)
	`, bookmark)
	if err != nil {
		return nil, errors.E(err, "cannot insert row")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.E(err, "cannot load last inserted ID")
	}
	err = b.db.Get(bookmark, `
		SELECT
			*
		FROM
			bookmarks
		WHERE
			id = $1
	`, id)
	if err != nil {
		return nil, errors.E(err, "cannot reload inserted row")
	}
	u, err := url.Parse(bookmark.URL)
	if err == nil {
		bookmark.Host = u.Host
	}
	return bookmark, errors.E(err, "cannot parse URL")
}

// GetByID loads one bookmark.
func (b *bookmarkDAO) GetByID(id int64) (*Bookmark, error) {
	bookmark := &Bookmark{}
	err := b.db.Get(bookmark, `
	SELECT
		*
	FROM
		bookmarks
	WHERE
		id = $1
	`, id)
	if err != nil {
		return nil, errors.E(err, "cannot find row")
	}
	u, err := url.Parse(bookmark.URL)
	if err == nil {
		bookmark.Host = u.Host
	}
	return bookmark, errors.E(err, "cannot parse URL")
}

// Update one bookmark.
func (b *bookmarkDAO) Update(bookmark *Bookmark) error {
	_, err := b.db.NamedExec(`
		UPDATE bookmarks
		SET
			url = :url,
			last_status_code = :last_status_code,
			last_status_check = :last_status_check,
			last_status_reason = :last_status_reason,
			title = :title,
			inbox = :inbox
		WHERE
			id = :id
	`, bookmark)
	return errors.E(err)
}

// Delete one bookmark.
func (b *bookmarkDAO) Delete(bookmark *Bookmark) error {
	_, err := b.db.NamedExec(`DELETE FROM bookmarks WHERE id = :id`, bookmark)
	return errors.E(err)
}
