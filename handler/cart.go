package handler

import (
	"context"
	"github.com/ShiyuCheng2018/cart/domain/model"
	"github.com/ShiyuCheng2018/cart/domain/service"
	cart "github.com/ShiyuCheng2018/cart/proto/cart"
	"github.com/ShiyuCheng2018/common"
)

type Cart struct {
	CartDataService service.ICartDataService
}

func (h *Cart) AddCart(ctx context.Context, request *cart.CartInfo, response *cart.ResponseAddCart) (err error) {
	cart := &model.Cart{}
	err = common.SwapTo(request, cart)
	if err != nil {
		return err
	}
	response.Id, err = h.CartDataService.AddCart(cart)
	return err
}

func (h *Cart) CleanCart(ctx context.Context, request *cart.Clean, response *cart.Response) (err error) {
	if err := h.CartDataService.CleanCart(request.UserId); err != nil {
		return nil
	}
	response.Meg = "Cart Cleaned"
	return nil
}

func (h *Cart) Increase(ctx context.Context, request *cart.Item, response *cart.Response) (err error) {
	if err := h.CartDataService.Increase(request.Id, request.ChangeNum); err != nil {
		return err
	}
	response.Meg = "Added item to cart successfully."
	return nil
}

func (h *Cart) Decrease(ctx context.Context, request *cart.Item, response *cart.Response) (err error) {
	if err := h.CartDataService.Decrease(request.Id, request.ChangeNum); err != nil {
		return err
	}
	response.Meg = "Removed item from cart successfully."
	return nil
}

func (h *Cart) DeleteItemById(ctx context.Context, request *cart.CartID, response *cart.Response) (err error) {
	if err := h.CartDataService.DeleteCart(request.Id); err != nil {
		return err
	}
	response.Meg = "Deleted item from cart successfully."
	return nil
}

func (h *Cart) GetAllCarts(ctx context.Context, request *cart.CartFindAll, response *cart.CartAll) error {
	cartAll, err := h.CartDataService.FindAllCart(request.UserId)
	if err != nil {
		return err
	}

	for _, v := range cartAll {
		cart := &cart.CartInfo{}
		if err := common.SwapTo(v, cart); err != nil {
			return err
		}

		response.CartInfo = append(response.CartInfo, cart)
	}
	return nil
}
