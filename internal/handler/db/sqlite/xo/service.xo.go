package xo

// Code generated by xo. DO NOT EDIT.

import (
	"context"
)

// Service represents a row from 'services'.
type Service struct {
	Parent      string `json:"parent"`      // parent
	Name        string `json:"name"`        // name
	Description string `json:"description"` // description
	UID         string `json:"uid"`         // uid
	Generation  int64  `json:"generation"`  // generation
	URI         string `json:"uri"`         // uri
	CreatedAt   Time   `json:"created_at"`  // created_at
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Service] exists in the database.
func (s *Service) Exists() bool {
	return s._exists
}

// Deleted returns true when the [Service] has been marked for deletion
// from the database.
func (s *Service) Deleted() bool {
	return s._deleted
}

// Insert inserts the [Service] to the database.
func (s *Service) Insert(ctx context.Context, db DB) error {
	switch {
	case s._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case s._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (manual)
	const sqlstr = `INSERT INTO services (` +
		`parent, name, description, uid, generation, uri, created_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`)`
	// run
	logf(sqlstr, s.Parent, s.Name, s.Description, s.UID, s.Generation, s.URI, s.CreatedAt)
	if _, err := db.ExecContext(ctx, sqlstr, s.Parent, s.Name, s.Description, s.UID, s.Generation, s.URI, s.CreatedAt); err != nil {
		return logerror(err)
	}
	// set exists
	s._exists = true
	return nil
}

// Update updates a [Service] in the database.
func (s *Service) Update(ctx context.Context, db DB) error {
	switch {
	case !s._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case s._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE services SET ` +
		`description = $1, uid = $2, generation = $3, uri = $4, created_at = $5 ` +
		`WHERE parent = $6 AND name = $7`
	// run
	logf(sqlstr, s.Description, s.UID, s.Generation, s.URI, s.CreatedAt, s.Parent, s.Name)
	if _, err := db.ExecContext(ctx, sqlstr, s.Description, s.UID, s.Generation, s.URI, s.CreatedAt, s.Parent, s.Name); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Service] to the database.
func (s *Service) Save(ctx context.Context, db DB) error {
	if s.Exists() {
		return s.Update(ctx, db)
	}
	return s.Insert(ctx, db)
}

// Upsert performs an upsert for [Service].
func (s *Service) Upsert(ctx context.Context, db DB) error {
	switch {
	case s._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO services (` +
		`parent, name, description, uid, generation, uri, created_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`)` +
		` ON CONFLICT (parent, name) DO ` +
		`UPDATE SET ` +
		`description = EXCLUDED.description, uid = EXCLUDED.uid, generation = EXCLUDED.generation, uri = EXCLUDED.uri, created_at = EXCLUDED.created_at `
	// run
	logf(sqlstr, s.Parent, s.Name, s.Description, s.UID, s.Generation, s.URI, s.CreatedAt)
	if _, err := db.ExecContext(ctx, sqlstr, s.Parent, s.Name, s.Description, s.UID, s.Generation, s.URI, s.CreatedAt); err != nil {
		return logerror(err)
	}
	// set exists
	s._exists = true
	return nil
}

// Delete deletes the [Service] from the database.
func (s *Service) Delete(ctx context.Context, db DB) error {
	switch {
	case !s._exists: // doesn't exist
		return nil
	case s._deleted: // deleted
		return nil
	}
	// delete with composite primary key
	const sqlstr = `DELETE FROM services ` +
		`WHERE parent = $1 AND name = $2`
	// run
	logf(sqlstr, s.Parent, s.Name)
	if _, err := db.ExecContext(ctx, sqlstr, s.Parent, s.Name); err != nil {
		return logerror(err)
	}
	// set deleted
	s._deleted = true
	return nil
}

// ServicesByCreatedAt retrieves a row from 'services' as a [Service].
//
// Generated from index 'created_at_desc'.
func ServicesByCreatedAt(ctx context.Context, db DB, createdAt Time) ([]*Service, error) {
	// query
	const sqlstr = `SELECT ` +
		`parent, name, description, uid, generation, uri, created_at ` +
		`FROM services ` +
		`WHERE created_at = $1`
	// run
	logf(sqlstr, createdAt)
	rows, err := db.QueryContext(ctx, sqlstr, createdAt)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Service
	for rows.Next() {
		s := Service{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&s.Parent, &s.Name, &s.Description, &s.UID, &s.Generation, &s.URI, &s.CreatedAt); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// ServiceByParentName retrieves a row from 'services' as a [Service].
//
// Generated from index 'sqlite_autoindex_services_1'.
func ServiceByParentName(ctx context.Context, db DB, parent, name string) (*Service, error) {
	// query
	const sqlstr = `SELECT ` +
		`parent, name, description, uid, generation, uri, created_at ` +
		`FROM services ` +
		`WHERE parent = $1 AND name = $2`
	// run
	logf(sqlstr, parent, name)
	s := Service{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, parent, name).Scan(&s.Parent, &s.Name, &s.Description, &s.UID, &s.Generation, &s.URI, &s.CreatedAt); err != nil {
		return nil, logerror(err)
	}
	return &s, nil
}
