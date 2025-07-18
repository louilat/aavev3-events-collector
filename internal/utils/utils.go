package utils

import (
	"math/big"
	"time"
)

func GetBlockReferences() map[int64]*big.Int {
	ref := map[int64]*big.Int{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(16532040),
		time.Date(2023, 2, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(16731850),
		time.Date(2023, 3, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(16952386),
		time.Date(2023, 4, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(17164064),
		time.Date(2023, 5, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(17384045),
		time.Date(2023, 6, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(17597293),
		time.Date(2023, 7, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(17818219),
		time.Date(2023, 8, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(18039775),
		time.Date(2023, 9, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(18253762),
		time.Date(2023, 10, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(18475328),
		time.Date(2023, 11, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(18689640),
		time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(18910671),

		time.Date(2024, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(19131664),
		time.Date(2024, 2, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(19338397),
		time.Date(2024, 3, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(19559071),
		time.Date(2024, 4, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(19773350),
		time.Date(2024, 5, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(19995041),
		time.Date(2024, 6, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(20209740),
		time.Date(2024, 7, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(20431767),
		time.Date(2024, 8, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(20653787),
		time.Date(2024, 9, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix():  big.NewInt(20868714),
		time.Date(2024, 10, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(21090866),
		time.Date(2024, 11, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(21305728),
		time.Date(2024, 12, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(21527686),

		time.Date(2025, 1, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(21749743),
		time.Date(2025, 2, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(21950084),
		time.Date(2025, 3, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(22172128),
		time.Date(2025, 4, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(22387076),
		time.Date(2025, 5, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(22607927),
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(22822463),
		time.Date(2025, 7, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)).Unix(): big.NewInt(22948418), // To change, month not finished yet
	}
	return ref
}
