package cli_test

import (
	"testing"

	util "github.com/Pylons-tech/pylons/testutil/cli"
	"github.com/Pylons-tech/pylons/testutil/network"
	"github.com/Pylons-tech/pylons/x/pylons/client/cli"
	"github.com/Pylons-tech/pylons/x/pylons/types"
	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
)

func settingPylonsAccountAndTrades(ctx client.Context, t *testing.T, net *network.Network, n int) string {
	acc := util.GenerateAddressesInKeyring(ctx.Keyring, 1)
	common := util.CommonArgs(acc[0].String(), net)
	username := "testUsername"
	token := "testToken"

	types.UpdateAppCheckFlagTest(types.FlagTrue)

	args := []string{}
	args = append(args, username, token, "")
	args = append(args, common...)

	_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateAccount(), args)
	require.NoError(t, err)

	coinInputs := "[{\"coins\": [{\"denom\": \"test\", \"amount\": \"1\"}]}]"
	coinOutputs := "[]"
	itemInputs := "[]"
	itemOutputs := "[]"
	extraInfo := "extra info"

	args = []string{}
	args = append(args, coinInputs)
	args = append(args, itemInputs)
	args = append(args, coinOutputs)
	args = append(args, itemOutputs)
	args = append(args, extraInfo)
	args = append(args, common...)

	for i := 0; i < n; i++ {
		_, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdCreateTrade(), args)
		require.NoError(t, err)
	}

	return acc[0].String()

}

func TestCmdListTradesByCreator(t *testing.T) {
	net := network.New(t)
	ctx := net.Validators[0].ClientCtx

	tradeCount := 3
	address := settingPylonsAccountAndTrades(ctx, t, net, tradeCount)

	for _, tc := range []struct {
		desc        string
		creator     string
		shouldErr   bool
		shouldFound bool
	}{
		{
			desc:        "Found",
			creator:     address,
			shouldErr:   false,
			shouldFound: true,
		},
		{
			desc:        "Not found",
			creator:     types.GenTestBech32FromString("Not found"),
			shouldErr:   false,
			shouldFound: false,
		},
		{
			desc:        "Invalid",
			creator:     "Invalid",
			shouldErr:   true,
			shouldFound: false,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{}
			args = append(args, tc.creator)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTradesByCreator(), args)
			if tc.shouldErr {
				require.Error(t, err)
				var resp types.QueryListTradesByCreatorResponse
				require.Error(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			} else if tc.shouldFound {
				require.NoError(t, err)
				var resp types.QueryListTradesByCreatorResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp)) // still error
				require.Equal(t, len(resp.Trades), tradeCount)

				for i := 0; i < tradeCount; i++ {
					require.Equal(t, resp.Trades[0].Creator, address)
				}
			} else {
				require.NoError(t, err)
				var resp types.QueryListTradesByCreatorResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp)) // still error
				require.Equal(t, len(resp.Trades), 0)
			}
		})
	}

}
