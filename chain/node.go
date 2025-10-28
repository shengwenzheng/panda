package chain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Node struct {
	ctx    context.Context
	client *ethclient.Client
}

func NewNode(ctx context.Context,rpcURL string)(*Node,error){
	client,err := ethclient.DialContext(ctx,rpcURL)
	if err != nil {
		return nil, fmt.Errorf("connect to node failed: %w", err) // “%w” wrapping 错误，使错误更加明晰
	}

	return &Node{ctx: ctx, client: client},nil
}

func (n *Node) GetBlockByNumber(number int64) (*types.Block,error) {
	block,err := n.client.BlockByNumber(n.ctx,nil)
	if err != nil {
		return nil, fmt.Errorf("get block failed: %w",err)
	}

	return block,nil
}

func (n *Node) GetLatestBlock() (*types.Block,error) {
	block,err := n.client.BlockByNumber(n.ctx,nil)
	if err != nil {
		return nil,err
	}

	return block, nil
}

func (n *Node) GetTxsByBlock(number int64) ([]*types.Transaction,error) {
	block,err := n.GetBlockByNumber(number)
	if err != nil {
		return nil, err
	}
	return block.Transactions(), nil
}

func (n *Node) GetTxByHash(txHash string) (*types.Transaction,bool,error) {
	hash := common.HexToHash(txHash)
	tx,isPending,err := n.client.TransactionByHash(n.ctx,hash)
	if err != nil {
		return nil,false,fmt.Errorf("get tx failed: %w",err)
	}

	return tx, isPending, nil
}

func (n *Node) GetChainID() (uint64,error) {
	chainID,err := n.client.ChainID(n.ctx)
	if err != nil {
		return 0, fmt.Errorf("get chainID failed %w",err)
	}

	return chainID.Uint64(),nil
}

func (n *Node) GetBlockHeight() (uint64,error) {
	block,err := n.GetLatestBlock()
	if err != nil {
		return 0,fmt.Errorf("get blockheight failed %w",err)
	}

	return block.NumberU64(),nil
}