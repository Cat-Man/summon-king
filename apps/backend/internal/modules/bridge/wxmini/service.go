package wxmini

import "fmt"

type Service struct{}

type LoginResult struct {
	UnifiedToken string `json:"unified_token"`
	Channel      string `json:"channel"`
}

type PayParams struct {
	OrderNo string `json:"order_no"`
	PaySign string `json:"pay_sign"`
	Channel string `json:"channel"`
}

func NewService() *Service { return &Service{} }

func (s *Service) ExchangeLogin(code string) LoginResult {
	return LoginResult{UnifiedToken: fmt.Sprintf("wxmini-token-%s", code), Channel: "wxmini"}
}

func (s *Service) BuildPayParams(orderNo string) PayParams {
	return PayParams{OrderNo: orderNo, PaySign: fmt.Sprintf("sign-%s", orderNo), Channel: "wxmini"}
}
