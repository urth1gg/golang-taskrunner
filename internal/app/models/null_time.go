package models

import (
	"time"
	"fmt"
)

type NullTime struct {
    Time  time.Time
    Valid bool // Valid is true if Time is not NULL
}

func (nt *NullTime) Scan(value interface{}) error {
    if value == nil {
        nt.Time, nt.Valid = time.Time{}, false
        return nil
    }
    nt.Valid = true
    switch v := value.(type) {
    case time.Time:
        nt.Time = v
    case []byte:
        // Assuming the date-time format in your database is "2006-01-02 15:04:05"
        t, err := time.Parse("2006-01-02 15:04:05", string(v))
        if err != nil {
            return err
        }
        nt.Time = t
    default:
        return fmt.Errorf("invalid Scan Source, unable to decode database value for created_at")
    }
    return nil
}