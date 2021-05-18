package encoder

import (
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
)

// EncodeDecimal
// Encodes a 256-bit integer as a number with mantissa and a decimal exponent.
// Exponent is 4 bits is packed in front of 12-bit mantissa.
// The original number can be recovered using following formula: V = M * 10^E
func EncodeDecimal(value models.Uint256) (uint16, error) {
	mantissa := new(big.Int).Set(&value.Int)
	exponent := big.NewInt(0)

	for i := 0; i < 15; i++ {
		if mantissa.Cmp(big.NewInt(0)) != 0 && big.NewInt(0).Mod(mantissa, big.NewInt(10)).Cmp(big.NewInt(0)) == 0 {
			mantissa.Div(mantissa, big.NewInt(10))
			exponent.Add(exponent, big.NewInt(1))
		} else {
			break
		}
	}

	if mantissa.Cmp(big.NewInt(0xfff)) > 0 {
		return 0, fmt.Errorf("value is not encodable as multi-precission decimal")
	}

	return uint16(exponent.Uint64())<<12 + uint16(mantissa.Uint64()), nil
}

// DecodeDecimal
// Decodes a 256-bit integer from a number with mantissa and a decimal exponent.
// Exponent is 4 bits is packed in front of 12-bit mantissa.
// The original number can be recovered using following formula: V = M * 10^E
func DecodeDecimal(value uint16) models.Uint256 {
	exponent := value >> 12
	mantissa := value & 0x0FFF // mantissa bitmask

	m := big.NewInt(int64(mantissa))
	exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(exponent)), nil)

	return models.MakeUint256FromBig(*m.Mul(m, exp))
}
