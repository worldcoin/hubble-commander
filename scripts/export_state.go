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

func ExportStateLeaves(filePath string) (err error) {
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

	leavesCount, err := exportStateLeaves(storage, file)
	if err != nil {
		return err
	}

	log.Infof("exported %d state leaves", leavesCount)
	return nil
}

func exportStateLeaves(storage *st.Storage, file *os.File) (int, error) {
	writer := bufio.NewWriter(file)
	_, err := writer.WriteString("[\n")
	if err != nil {
		return 0, err
	}
	count := 0

	err = storage.StateTree.IterateLeaves(func(stateLeaf *models.StateLeaf) error {
		if count > 0 {
			err = writer.WriteByte(',')
			if err != nil {
				return err
			}
		}
		count++

		return writeLeaf(writer, stateLeaf)
	})
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

func writeLeaf(writer *bufio.Writer, leaf *models.StateLeaf) error {
	bytes, err := json.MarshalIndent(leaf, "", "\t")
	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	return err
}
