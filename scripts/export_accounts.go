package scripts

import (
	"bufio"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ExportAccounts(filePath string) error {
	return exportLeaves(filePath, exportAndCountAccounts)
}

func exportAndCountAccounts(storage *st.Storage, writer *bufio.Writer) (int, error) {
	count := 0
	err := storage.AccountTree.IterateLeaves(func(accountLeaf *models.AccountLeaf) error {
		if count > 0 {
			err := writer.WriteByte(',')
			if err != nil {
				return err
			}
		}
		count++

		return writeData(writer, accountLeaf)
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
