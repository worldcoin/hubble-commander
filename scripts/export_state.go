package scripts

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
)

type exportFunc func(storage *st.Storage, writer *bufio.Writer) (int, error)

func ExportStateLeaves(filePath string) error {
	return exportLeaves(filePath, exportAndCountStateLeaves)
}

func exportLeaves(filePath string, exportDataFunc exportFunc) (err error) {
	cfg := config.GetCommanderConfigAndSetupLogger()
	storage, err := st.NewStorage(cfg)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := file.Close()
		if err != nil {
			err = closeErr
		}
	}()

	leavesCount, err := exportData(storage, file, exportDataFunc)
	if err != nil {
		return err
	}

	log.Infof("exported %d leaves", leavesCount)
	return nil
}

func exportData(storage *st.Storage, file *os.File, exportDataFunc exportFunc) (int, error) {
	writer := bufio.NewWriter(file)
	err := writer.WriteByte('[')
	if err != nil {
		return 0, err
	}

	count, err := exportDataFunc(storage, writer)
	if err != nil {
		return 0, err
	}

	err = writer.WriteByte(']')
	if err != nil {
		return 0, err
	}

	err = writer.Flush()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func exportAndCountStateLeaves(storage *st.Storage, writer *bufio.Writer) (int, error) {
	count := 0
	err := storage.StateTree.IterateLeaves(func(stateLeaf *models.StateLeaf) error {
		if count > 0 {
			err := writer.WriteByte(',')
			if err != nil {
				return err
			}
		}
		count++

		return writeData(writer, stateLeaf)
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func writeData(writer *bufio.Writer, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	return err
}
