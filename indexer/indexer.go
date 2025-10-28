package indexer

import (
	"context"
	"time"
	"log/slog"

	"gorm.io/gorm"
	"panda.com/chain"
	"panda.com/database/orm"
)

const safeConfirmNum = 20

type Indexer struct {
	ctx      context.Context
	interval int64
	node     *chain.Node
	db       *gorm.DB
}

func NewIndexer(ctx context.Context,interval int64,node *chain.Node,db *gorm.DB) *Indexer {
	return &Indexer {
		ctx:      ctx,
		interval: interval,
		node:     node,
		db:       db,
	}
}

func (i *Indexer) Start() {
	ticker := time.NewTicker(time.Duration(i.interval) * time.Second)
	defer ticker.Stop()

	for {
		select{
		case <-ticker.C:
			if err := i.sync(); err != nil {
				slog.Error("indexer sync error","error",err)
			}
		case <-i.ctx.Done():
			slog.Info("indexer context done, exiting...")
			return
		}
	}
}

func (i *Indexer) sync() error {
	latestheight,err := i.localHeight()
	if err != nil {
		return err
	}
	headHeight,err := i.node.GetBlockHeight()
	if err != nil {
		return  err
	}
	headHeight -= safeConfirmNum
	if latestheight >= int64(headHeight) {
		return nil
	}

	return nil
}

func (i *Indexer) localHeight() (int64,error) {
	var latestChain orm.Chain
	if err := i.db.First(&latestChain).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		if err := i.db.Create(&orm.Chain{
			Height: 0,
		}).Error; err != nil {
			return 0, err
		}

		return 0, nil
	}

	return latestChain.Height, nil
}