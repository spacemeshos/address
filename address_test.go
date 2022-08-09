package address_test

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"

	"github.com/spacemeshos/address"
)

func init() {
	address.DefaultTestAddressConfig()
}

func TestAddress_NewAddress(t *testing.T) {
	t.Parallel()
	table := []struct {
		name string
		src  string
		err  error
	}{
		{
			name: "correct",
			src:  "stest1qqqqqqy0dd83jemjmfj3ghjndxm0ndh0z2rymaqyp0gu3",
		},
		{
			name: "corect address, but no reserved space",
			src:  "stest1fejq2x3d79ukpkw06t7h6lndjuwzxdnj59npghsg43mh4",
			err:  address.ErrMissingReservedSpace,
		},
		{
			name: "incorrect cahrset", // chars iob is not supported by bech32.
			src:  "stest1fejq2x3d79ukpkw06t7h6lniobuwzxdnj59nphsywyusj",
			err:  address.ErrDecodeBech32,
		},
		{
			name: "missing hrp",
			src:  "1fejq2x3d79ukpkw06t7h6lniobuwzxdnj59nphsywyusj",
			err:  address.ErrDecodeBech32,
		},
		{
			name: "to big length",
			src:  "stest1qw508d6e1qejxtdg4y5r3zarvax8wucu",
			err:  address.ErrWrongAddressLength,
		},
		{
			name: "to small length",
			src:  "stest1qw504y5r3zarva2v6vda",
			err:  address.ErrWrongAddressLength,
		},
		{
			name: "wrong network",
			src:  "sut1fejq2x3d79ukpkw06t7h6lndjuwzxdnj59npghsldfkky",
			err:  address.ErrUnsupportedNetwork,
		},
	}

	for _, testCase := range table {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := address.StringToAddress(testCase.src)
			if testCase.err != nil {
				require.ErrorContains(t, err, testCase.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAddress_SetBytes(t *testing.T) {
	t.Parallel()
	t.Run("generate from correct string", func(t *testing.T) {
		srcAddr, err := address.StringToAddress("stest1qqqqqqrs60l66w5uksxzmaznwq6xnhqfv56c28qlkm4a5")
		require.NoError(t, err)
		newAddr := address.GenerateAddress(srcAddr.Bytes())
		require.NoError(t, err)
		checkAddressesEqual(t, srcAddr, newAddr)
	})
	t.Run("generate from small byte slice", func(t *testing.T) {
		newAddr := address.GenerateAddress([]byte{1})
		require.NotEmpty(t, newAddr.String())
	})
	t.Run("generate from large byte slice length", func(t *testing.T) {
		newAddr := address.GenerateAddress([]byte("12345678901234567890123456789012345678901234567890"))
		require.NotEmpty(t, newAddr.String())
	})
	t.Run("generate from correct byte slice", func(t *testing.T) {
		data := make([]byte, address.AddressLength)
		for i := range data {
			data[i] = byte(i)
		}
		newAddr := address.GenerateAddress(data)
		srcAddr, err := address.StringToAddress(newAddr.String())
		require.NoError(t, err)
		checkAddressesEqual(t, newAddr, srcAddr)
	})
}

func TestAddress_GenerateAddress(t *testing.T) {
	t.Parallel()
	t.Run("generate from public key", func(t *testing.T) {
		srcAddr := address.GenerateAddress(RandomBytes(32))
		newAddr, err := address.StringToAddress(srcAddr.String())
		require.NoError(t, err)
		checkAddressesEqual(t, srcAddr, newAddr)
	})
	t.Run("generate from string", func(t *testing.T) {
		srcAddr := address.GenerateAddress([]byte("some random very very long string"))
		newAddr, err := address.StringToAddress(srcAddr.String())
		require.NoError(t, err)
		checkAddressesEqual(t, srcAddr, newAddr)
	})
}

func TestAddress_ReservedBytesOnTop(t *testing.T) {
	for i := 0; i < 100; i++ {
		addr := address.GenerateAddress(RandomBytes(i))
		for j := 0; j < address.AddressReservedSpace; j++ {
			require.Zero(t, addr.Bytes()[j])
		}
	}
}

func checkAddressesEqual(t *testing.T, addrA, addrB address.Address) {
	require.Equal(t, addrA.Bytes(), addrB.Bytes())
	require.Equal(t, addrA.String(), addrB.String())
	require.Equal(t, addrA.IsEmpty(), addrB.IsEmpty())
}

// RandomBytes generates random data in bytes for testing.
func RandomBytes(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}
