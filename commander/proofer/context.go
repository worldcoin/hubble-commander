package proofer

import (
	st "github.com/Worldcoin/hubble-commander/storage"
)

type Context struct {
	storage *st.Storage
}

func NewContext(storage *st.Storage) *Context {
	return &Context{storage: storage}
}
