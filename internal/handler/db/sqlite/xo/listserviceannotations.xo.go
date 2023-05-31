package xo

// Code generated by xo. DO NOT EDIT.

import (
	"context"
)

// ListServiceAnnotations represents a row from 'list_service_annotations'.
type ListServiceAnnotations struct {
	ServiceParent string `json:"service_parent"` // service_parent
	ServiceName   string `json:"service_name"`   // service_name
	Key           string `json:"key"`            // key
	Value         string `json:"value"`          // value
}

// ListServiceAnnotationsByParentName runs a custom query, returning results as [ListServiceAnnotations].
func ListServiceAnnotationsByParentName(ctx context.Context, db DB, parent, name string) ([]*ListServiceAnnotations, error) {
	// query
	const sqlstr = `SELECT` +
		`    service_parent,` +
		`    service_name,` +
		`    key,` +
		`    value` +
		`  FROM service_annotations` +
		`  WHERE` +
		`    service_parent = $1 AND` +
		`    service_name = $2`
	// run
	logf(sqlstr, parent, name)
	rows, err := db.QueryContext(ctx, sqlstr, parent, name)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*ListServiceAnnotations
	for rows.Next() {
		var lsa ListServiceAnnotations
		// scan
		if err := rows.Scan(&lsa.ServiceParent, &lsa.ServiceName, &lsa.Key, &lsa.Value); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &lsa)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}
