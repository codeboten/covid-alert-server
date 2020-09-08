package persistence

import "database/sql"
import "time"

// Event the event that we are to log
type Event struct {
	identifier EventType
	deviceType DeviceType
	date       time.Time
	count      int
}

func (c *conn) SaveEvent(event Event) error {

	if err := saveEvent(c.db, event); err != nil {
		return err
	}
	return nil
}

// DeviceType the type of the device the event was generated by
type DeviceType string

const (
	android DeviceType = "Android"
	apple   DeviceType = "iOS"
	server  DeviceType = "server"
)

// EventType the type of the event that happened
type EventType string

const (
	keyClaimed   EventType = "keyClaimed"
	keyGenerated EventType = "keyGenerated"
)

func saveEvent(db *sql.DB, e Event) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO events
		(identifier, device_type, date, count)
		VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE count = count + ?`,
		e.identifier, e.deviceType, e.date.Format("2006-01-02"), e.count, e.count); err != nil {

		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
