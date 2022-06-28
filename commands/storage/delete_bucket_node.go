package commands

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
	"strconv"
)

type DeleteBucketNodeRecursivelyCommand struct {
	Node     Node
	Path     string
	Database Database
}

func (c DeleteBucketNodeRecursivelyCommand) Run() ICommandError {
	nodes, err := c.Database.GetNodesFromNode(c.Node.Id)
	if err != nil && err != sql.ErrNoRows {
		return NewError(http.StatusInternalServerError, err)
	}

	for _, node := range nodes {
		var err ICommandError

		path := filepath.Join(c.Path, node.Filename)

		switch node.Filetype {
		case "directory":
			err = DeleteBucketNodeRecursivelyCommand{
				Node:     node,
				Path:     path,
				Database: c.Database,
			}.Run()
		default:
			err = NewTransaction([]Command{
				DeleteBucketNodeCommand{
					Node:     node,
					Database: c.Database,
				},
				DeleteBucketNodeInFileSystemCommand{
					Node:     node,
					Path:     path,
					Database: c.Database,
				},
			}).Try()
		}

		if err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}
	}

	transactionError := NewTransaction([]Command{
		DeleteBucketNodeCommand{
			Node:     c.Node,
			Database: c.Database,
		},
		DeleteBucketNodeInFileSystemCommand{
			Node:     c.Node,
			Path:     c.Path,
			Database: c.Database,
		},
	}).Try()

	if transactionError != nil {
		return NewError(http.StatusInternalServerError, transactionError.Error())
	}
	return nil
}

func (c DeleteBucketNodeRecursivelyCommand) Revert() ICommandError {
	// TODO: Revert file deletion
	return nil
}

type DeleteBucketNodeCommand struct {
	Node     Node
	Database Database
}

func (c DeleteBucketNodeCommand) Run() ICommandError {
	request := `
		BEGIN TRANSACTION;
		DELETE FROM buckets_nodes WHERE id = ?;
		DELETE FROM buckets_nodes_associations WHERE to_node = ?;
		COMMIT TRANSACTION;
	`

	_, err := c.Database.Instance.Exec(request, c.Node.Id, c.Node.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c DeleteBucketNodeCommand) Revert() ICommandError {
	// TODO: Revert file deletion
	return nil
}

type DeleteBucketNodeInFileSystemCommand struct {
	Node     Node
	Path     string
	Database Database
}

func (c DeleteBucketNodeInFileSystemCommand) Run() ICommandError {
	if len(c.Path) > 0 && c.Path[0] == '/' {
		c.Path = c.Path[1:]
	}

	c.Path = filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(c.Node.BucketId), c.Path)
	err := os.RemoveAll(c.Path)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c DeleteBucketNodeInFileSystemCommand) Revert() ICommandError {
	// TODO: Revert file deletion
	return nil
}