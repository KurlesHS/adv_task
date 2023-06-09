package postgres

import (
	"context"
	"database/sql"
	"fmt"
	errormessage "kurles/adv_task/pkg/error_message"
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
	if err != nil {
		err = errormessage.NewError(errormessage.Db, err.Error())
	}
	return
}

func (repo *PgAdvertRepo) GetAdverts(ctx context.Context, page int, sortBy model.SortBy, sortOrder model.SortOrder) (res []model.Advert, err error) {
	if page < 1 {
		err = errormessage.NewError(errormessage.BadRequest, fmt.Sprintf("wrong page: %v", page))
		return
	}

	var sortByV string
	if sortBy == model.Date {
		sortByV = "create_timestamp"
	} else {
		sortByV = "price"
	}

	sortOrderV := ""
	if sortOrder == model.Desc {
		sortOrderV = "desc"
	}

	q := fmt.Sprintf("select id, title, price, (select link from photos where adv_id = id limit 1) "+
		"from adverts ad "+
		"order by %v %v limit %v offset %v;", sortByV, sortOrderV, pageSize, pageSize*(page-1))
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		err = errormessage.NewError(errormessage.Db, err.Error())
		return
	}
	defer rows.Close()
	res = make([]model.Advert, 0)
	for rows.Next() {
		var record model.Advert
		var title, link sql.NullString
		err = rows.Scan(&record.Id, &title, &record.Price, &link)
		if err != nil {
			err = errormessage.NewError(errormessage.Internal, err.Error())
			return
		}
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
		err = errormessage.NewError(errormessage.Db, err.Error())
		return
	}
	defer rows.Close()
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
			err = errormessage.NewError(errormessage.Db, err.Error())
			return
		}

		for rows.Next() {
			err = rows.Scan(&link)
			if err == nil {
				adv.Photos = append(adv.Photos, link.String)
			}
		}
		rows.Close()
	} else {
		err = errormessage.NewError(errormessage.NotFound, fmt.Sprintf("advert with id %v is not found", advId))
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
		err = errormessage.NewError(errormessage.Db, err.Error())
		_ = tx.Rollback()
		return
	}
	for _, v := range advert.Photos {
		_, err = tx.ExecContext(ctx, "insert into photos (adv_id, link) values ($1, $2);", id, v)
		if err != nil {
			err = errormessage.NewError(errormessage.Db, err.Error())
			_ = tx.Rollback()
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		err = errormessage.NewError(errormessage.Db, err.Error())
	}
	return
}

func (repo *PgAdvertRepo) ClearAllAdverts(ctx context.Context) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		err = errormessage.NewError(errormessage.Db, err.Error())
		return err
	}
	_, err = tx.ExecContext(ctx, "delete from adverts;")
	if err != nil {
		_ = tx.Rollback()
		err = errormessage.NewError(errormessage.Db, err.Error())
		return err
	}
	err = tx.Commit()
	if err != nil {
		err = errormessage.NewError(errormessage.Db, err.Error())
	}
	return err
}
