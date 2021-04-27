package spanner

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
)

func writeUsingDML(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT user (id, name) VALUES
                                (2, 'Melissa'),
                                (3, 'Russell'),
                                (4, 'Jacqueline'),
                                (5, 'Dylan')`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return err
	})
	return err
}
