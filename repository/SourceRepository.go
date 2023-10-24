package repository

import "database/sql"

type SourceRepository struct {
    db *sql.DB
}

type Source struct {
    Id int `field:"id"`
    Url string `field:"url"`
    Prefix string `field:"prefix"`
    Selector string `field:"selector"`
    Name string `field:"name"`
}

func NewSourceRepository(db *sql.DB) *SourceRepository {
    return &SourceRepository{
        db,
    }
}

func (repo SourceRepository) GetSources() ([]*Source, error) {
    sources := make([]*Source, 0)

    rows, err := repo.db.Query("SELECT id, url, prefix, selector, name FROM sources")
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        source := new(Source)
        err = rows.Scan(&source.Id, &source.Url, &source.Prefix, &source.Selector, &source.Name)
        if err != nil {
            return nil, err
        }
        sources = append(sources, source)
    }

    return sources, nil
}
