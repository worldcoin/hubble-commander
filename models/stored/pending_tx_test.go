package stored

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// nolint:gocritic
func (t PendingTx) Generate(rand *rand.Rand, size int) reflect.Value {
	hashBytes := make([]byte, 32)
	rand.Read(hashBytes)
	hash := common.BytesToHash(hashBytes)

	var receiveTime *models.Timestamp
	if rand.Intn(2) == 0 {
		randomUnixTime := rand.Int63()
		// This must be UTC() or else the equality check will fail,
		// see models/timestamp.go:Bytes()
		randomTime := time.Unix(randomUnixTime, 0).UTC()
		receiveTime = models.NewTimestamp(randomTime)
	}

	return reflect.ValueOf(PendingTx{
		Hash: hash,

		TxType:      txtype.Transfer,
		FromStateID: rand.Uint32(),
		Amount:      models.MakeUint256(1000),
		Fee:         models.MakeUint256(100),
		Nonce:       models.MakeUint256(0),
		Signature:   models.MakeRandomSignature(),
		ReceiveTime: receiveTime,

		Body: &TxTransferBody{
			ToStateID: rand.Uint32(),
		},
	})
}

// BatchedTx.SetBytes and FailedTx.SetBytes rely on this property!

func TestPendingTx_BytesLenMatchesBytes(t *testing.T) {
	f := func() bool {
		valuePendingTx, ok := quick.Value(
			reflect.TypeOf(PendingTx{}),
			// nolint:gosec
			rand.New(rand.NewSource(time.Now().Unix())),
		)
		require.True(t, ok)
		require.NotNil(t, valuePendingTx)
		pendingTx := valuePendingTx.Interface().(PendingTx)

		bytesLen := pendingTx.BytesLen()
		bytes := pendingTx.Bytes()
		require.Equal(t, len(bytes), bytesLen)

		var deserialized PendingTx
		err := deserialized.SetBytes(bytes)
		require.Nil(t, err)
		require.Equal(t, pendingTx, deserialized)

		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
