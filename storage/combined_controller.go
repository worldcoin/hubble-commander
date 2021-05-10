package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
)

type CombinedController struct {
	postgresTx *db.TxController
	badgerTx   *db.TxController
}

func NewCombinedController(postgresTx, badgerTx *db.TxController) *CombinedController {
	return &CombinedController{
		postgresTx: postgresTx,
		badgerTx:   badgerTx,
	}
}

func (c *CombinedController) Rollback() error {
	var err error
	c.postgresTx.Rollback(&err)
	c.badgerTx.Rollback(&err)
	return nil
}

func (c *CombinedController) Commit() error {
	if err := c.postgresTx.Commit(); err != nil {
		return err
	}
	return c.badgerTx.Commit()
}
