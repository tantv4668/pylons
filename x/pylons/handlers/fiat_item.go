package handlers

import (
	"github.com/MikeSofaer/pylons/x/pylons/keep"
	"github.com/MikeSofaer/pylons/x/pylons/msgs"
	"github.com/MikeSofaer/pylons/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FiatItemResponse struct {
	ItemID string `json:"ItemID"`
}

// HandlerMsgFiatItem is used to create item within 1 block execution
func HandlerMsgFiatItem(ctx sdk.Context, keeper keep.Keeper, msg msgs.MsgFiatItem) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}
	cook, err2 := keeper.GetCookbook(ctx, msg.CookbookID)
	if err2 != nil {
		return errInternal(err2)
	}
	if !cook.Sender.Equals(msg.Sender) {
		// TODO: should enable it if fiat_item should only be signed by game dev
		// return sdk.ErrUnauthorized("cookbook not owned by the sender").Result()
	}

	item := types.NewItem(msg.CookbookID, msg.Doubles, msg.Longs, msg.Strings, msg.Sender)

	if err := keeper.SetItem(ctx, *item); err != nil {
		return errInternal(err)
	}

	return marshalJson(FiatItemResponse{
		item.ID,
	})
}
