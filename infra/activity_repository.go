package infra

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type sqlActivity struct {
	Id        string    `db:"id"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	Name      string    `db:"name"`
	Content   string    `db:"content"`
	AuthorUrl string    `db:"author_url"`
	Raw       string    `db:"raw"`
}

type DefaultActivityRepo struct {
	db *sqlx.DB
}

func NewActivityRepo(db Database) repository.ActivityRepository {
	sqlxdb := db.Get()

	sqlxdb.MustExec(`
        CREATE TABLE IF NOT EXISTS activities (
			id TEXT PRIMARY KEY,
            type TEXT NOT NULL,
            created_at DATETIME NOT NULL,
            name TEXT,
            content TEXT NOT NULL,
            author_url TEXT NOT NULL,
            raw TEXT NOT NULL
        );
    `)

	return &DefaultActivityRepo{
		db: sqlxdb,
	}
}

func (repo *DefaultActivityRepo) Upsert(act *model.Activity) error {
	_, err := repo.db.NamedExec(`
        INSERT INTO activities 
        (id, type, created_at, name, content, author_url, raw)
        VALUES 
        (:id, :type, :created_at, :name, :content, :author_url, :raw)
        ON CONFLICT(id) DO UPDATE SET
            type=excluded.type,
            name=excluded.name,
            content=excluded.content,
            author_url=excluded.author_url,
            raw=excluded.raw
    `, sqlActivity{
		Id:        act.Id,
		Type:      act.Type,
		CreatedAt: act.CreatedAt,
		Name:      act.Name,
		Content:   act.Content,
		AuthorUrl: act.AuthorUrl,
		Raw:       act.Raw,
	})
	return err
}

func (repo *DefaultActivityRepo) ListRecent(page int, size int) ([]model.Activity, error) {
	var acts []sqlActivity
	err := repo.db.Select(&acts, `
        SELECT * 
        FROM activities
        ORDER BY created_at DESC
        LIMIT ?
        OFFSET ?
    `, size, page*size)
	result := make([]model.Activity, 0, len(acts))
	for _, a := range acts {
		act, err := repo.sqlActivityToActivity(a)
		if err != nil {
			return nil, err
		}
		result = append(result, act)
	}
	return result, err
}

func (r *DefaultActivityRepo) sqlActivityToActivity(a sqlActivity) (model.Activity, error) {
	return model.Activity{
		Id:        a.Id,
		Type:      a.Type,
		CreatedAt: a.CreatedAt,
		Name:      a.Name,
		Content:   a.Content,
		AuthorUrl: a.AuthorUrl,
		Raw:       a.Raw,
	}, nil
}
