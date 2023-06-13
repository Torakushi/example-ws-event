package persistence

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"time"

	"github.com/Torakushi/example-ws-events/internal/models"
)

func NewConnection(pgURL string) (*pg.DB, error) {
	opt, err := pg.ParseURL(pgURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL URL: %v", err)
	}

	db := pg.Connect(opt)
	if db == nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL")
	}

	return db, nil
}

type EventRepository struct {
	db *pg.DB
}

func NewEventRepository(db *pg.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) SaveEvent(event *models.Event) error {
	event.Timestamp = time.Now()

	res, err := r.db.Conn().Model(event).Insert()
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("no row inserted")
	}

	return nil
}

func (r *EventRepository) GetLastEvents(n int) ([]*models.Event, error) {
	var events []*models.Event

	err := r.db.Conn().Model(&events).Order("timestamp DESC").Limit(n).Select()
	if err != nil {
		return nil, err
	}

	return events, nil
}
