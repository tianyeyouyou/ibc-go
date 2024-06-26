package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

func TestValidateIBCDenom(t *testing.T) {
	testCases := []struct {
		name     string
		denom    string
		expError bool
	}{
		{"denom with trace hash", "ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", false},
		{"base denom", "uatom", false},
		{"base denom ending with '/'", "uatom/", false},
		{"base denom with single '/'s", "gamm/pool/1", false},
		{"base denom with double '/'s", "gamm//pool//1", false},
		{"non-ibc prefix with hash", "notibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", false},
		{"empty denom", "", true},
		{"denom 'ibc'", "ibc", true},
		{"denom 'ibc/'", "ibc/", true},
		{"invalid hash", "ibc/!@#$!@#", true},
	}

	for _, tc := range testCases {
		tc := tc

		err := types.ValidateIBCDenom(tc.denom)
		if tc.expError {
			require.Error(t, err, tc.name)
			continue
		}
		require.NoError(t, err, tc.name)
	}
}

func TestExtractDenomFromFullPath(t *testing.T) {
	testCases := []struct {
		name     string
		fullPath string
		expDenom types.Denom
	}{
		{"empty denom", "", types.Denom{}},
		{"base denom no slashes", "atom", types.Denom{Base: "atom"}},
		{"base denom with trailing slash", "atom/", types.Denom{Base: "atom/"}},
		{"base denom multiple trailing slash", "foo///bar//baz/atom/", types.Denom{Base: "foo///bar//baz/atom/"}},
		{"ibc denom one hop", "transfer/channel-0/atom", types.Denom{Base: "atom", Trace: []types.Trace{types.NewTrace("transfer", "channel-0")}}},
		{"ibc denom one hop trailing slash", "transfer/channel-0/atom/", types.Denom{Base: "atom/", Trace: []types.Trace{types.NewTrace("transfer", "channel-0")}}},
		{"ibc denom one hop multiple slashes", "transfer/channel-0//at/om/", types.Denom{Base: "/at/om/", Trace: []types.Trace{types.NewTrace("transfer", "channel-0")}}},
		{"ibc denom two hops", "transfer/channel-0/transfer/channel-60/atom", types.Denom{Base: "atom", Trace: []types.Trace{types.NewTrace("transfer", "channel-0"), types.NewTrace("transfer", "channel-60")}}},
		{"ibc denom two hops trailing slash", "transfer/channel-0/transfer/channel-60/atom/", types.Denom{Base: "atom/", Trace: []types.Trace{types.NewTrace("transfer", "channel-0"), types.NewTrace("transfer", "channel-60")}}},
		{"empty prefix", "/uatom", types.Denom{Base: "/uatom"}},
		{"empty identifiers", "//uatom", types.Denom{Base: "//uatom"}},
		{"base denom with single '/'", "erc20/0x85bcBCd7e79Ec36f4fBBDc54F90C643d921151AA", types.Denom{Base: "erc20/0x85bcBCd7e79Ec36f4fBBDc54F90C643d921151AA"}},
		{"trace info and base denom with single '/'", "transfer/channel-1/erc20/0x85bcBCd7e79Ec36f4fBBDc54F90C643d921151AA", types.Denom{Base: "erc20/0x85bcBCd7e79Ec36f4fBBDc54F90C643d921151AA", Trace: []types.Trace{types.NewTrace("transfer", "channel-1")}}},
		{"single trace identifier", "transfer/", types.Denom{Base: "transfer/"}},
		{"trace info with custom port", "customtransfer/channel-1/uatom", types.Denom{Base: "uatom", Trace: []types.Trace{types.NewTrace("customtransfer", "channel-1")}}},
		{"invalid path (1)", "channel-1/transfer/uatom", types.Denom{Base: "channel-1/transfer/uatom"}},
		{"invalid path (2)", "transfer/channel-1", types.Denom{Base: "transfer/channel-1"}},
		{"invalid path (3)", "transfer/channel-1/transfer/channel-2", types.Denom{Trace: []types.Trace{types.NewTrace("transfer", "channel-1"), types.NewTrace("transfer", "channel-2")}}},
		{"invalid path (4)", "transfer/channelToA/uatom", types.Denom{Base: "transfer/channelToA/uatom"}},
	}

	for _, tc := range testCases {
		tc := tc

		denom := types.ExtractDenomFromFullPath(tc.fullPath)
		require.Equal(t, tc.expDenom, denom, tc.name)
	}
}
