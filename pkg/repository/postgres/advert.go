package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"kurles/adv_task/pkg/model"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PgAdvertRepo struct {
	db *sqlx.DB
}

const (
	pageSize = 10
)

func New(host string, port int, dbName string, user string, pass string) (repo PgAdvertRepo, err error) {
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName)
	repo.db, err = sqlx.Open("postgres", conStr)
	return
}

func (repo *PgAdvertRepo) GetAdverts(ctx context.Context, page int, sortBy model.SortBy, desc bool) (res []model.Advert, err error) {
	if page < 1 {
		err = fmt.Errorf("wrong page: %v", page)
		return
	}
	q := fmt.Sprintf("select id, title, price, (select link from photos where adv_id = id limit 1) "+
		"from adverts ad "+
		"order by create_timestamp asc limit %v offset %v;", pageSize, pageSize*(page-1))
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return
	}
	for rows.Next() {
		var record model.Advert
		var title, link sql.NullString
		err = rows.Scan(&record.Id, &title, &record.Price, &link)
		if title.Valid {
			record.Title = title.String
		}
		if link.Valid {
			record.MainPhoto = link.String
		}
		res = append(res, record)
	}
	return
}

func (repo *PgAdvertRepo) GetAdvert(ctx context.Context, advId int64) (adv model.DetailedAdvert, err error) {
	q := "select title, description, price from adverts where id = $1;"
	rows, err := repo.db.QueryContext(ctx, q, advId)
	if err != nil {
		return
	}
	if rows.Next() {
		var title, description, link sql.NullString
		err = rows.Scan(&title, &description, &adv.Price)
		adv.Id = advId
		adv.Title = title.String
		adv.Description = description.String
		if err != nil {
			return
		}
		rows, err = repo.db.QueryContext(ctx, "select link from photos where adv_id = $1 order by photo_order;", advId)
		if err != nil {
			return
		}
		for rows.Next() {
			rows.Scan(&link)
			adv.Photos = append(adv.Photos, link.String)
		}
	}
	return
}

func (repo *PgAdvertRepo) InsertAdvert(ctx context.Context, advert model.DetailedAdvert) (id int64, err error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	q := "insert into adverts (title, description, price, create_timestamp) values ($1, $2, $3, $4) returning id;"
	row := tx.QueryRowContext(ctx, q, advert.Title, advert.Description, advert.Price, time.Now())
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return
	}
	for _, v := range advert.Photos {
		_, err = tx.ExecContext(ctx, "insert into photos (adv_id, link) values ($1, $2);", id, v)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit()
	return
}
