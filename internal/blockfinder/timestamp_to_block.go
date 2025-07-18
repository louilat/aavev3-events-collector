package blockfinder

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FindClosestBlocks(backend *ethclient.Client, tmstp uint64, ref *big.Int, step int64) (*big.Int, *big.Int, error) {
	ub := big.NewInt(0)
	lb := big.NewInt(0)
	ub.Set(ref)
	lb.Add(ref, big.NewInt(-step))

	lh, err := backend.HeaderByNumber(context.Background(), lb)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), err
	}

	uh, err := backend.HeaderByNumber(context.Background(), ub)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), err
	}

	for lh.Time > tmstp {
		// fmt.Printf("Lower block is %v\n", lb)
		// fmt.Printf("Upper block is %v\n", ub)
		// fmt.Println("------------------------------")
		ub.Set(lb)
		lb.Add(lb, big.NewInt(-step))

		lh, err = backend.HeaderByNumber(context.Background(), lb)
		if err != nil {
			return big.NewInt(0), big.NewInt(0), err
		}

		uh, err = backend.HeaderByNumber(context.Background(), ub)
		if err != nil {
			return big.NewInt(0), big.NewInt(0), err
		}
	}

	if uh.Time < tmstp {
		return big.NewInt(0), big.NewInt(0), errors.New("ref probably lower than target block")
	}

	mb := big.NewInt(0)
	var mh *types.Header
	s := big.NewInt(0)

	dif := big.NewInt(0)
	dif.Sub(ub, lb)

	for dif.Cmp(big.NewInt(1)) > 0 {
		// fmt.Printf("diff is %v\n", dif.String())
		s.Add(lb, ub)
		mb.Div(s, big.NewInt(2))
		// fmt.Printf("mb is %v\n", mb.String())

		mh, err = backend.HeaderByNumber(context.Background(), mb)
		if err != nil {
			return big.NewInt(0), big.NewInt(0), err
		}

		if mh.Time < tmstp {
			lb.Set(mb)
		} else {
			ub.Set(mb)
		}

		dif.Sub(ub, lb)
		if dif.Cmp(big.NewInt(0)) < 0 {
			return big.NewInt(0), big.NewInt(0), errors.New("failed dichotomy")
		}

	}

	return lb, ub, nil
}
