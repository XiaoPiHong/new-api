package controller

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

func paymentReturnPath(suffix string) string {
	base := strings.TrimRight(system_setting.ServerAddress, "/")
	return base + common.ThemeAwarePath(suffix)
}

// topUpPaymentReturnURL 生成充值完成后的浏览器返回地址。
// xph_order_no 使用 NewAPI 订单号，供自有前端查询最终状态；它不代表支付已经成功。
func topUpPaymentReturnURL(tradeNo string) (*url.URL, error) {
	returnAddress := strings.TrimSpace(operation_setting.GetPaymentSetting().TopUpReturnURL)
	if returnAddress == "" {
		returnAddress = paymentReturnPath("/console/log")
	}

	returnURL, err := url.Parse(returnAddress)
	if err != nil || returnURL.Host == "" || (returnURL.Scheme != "http" && returnURL.Scheme != "https") {
		return nil, fmt.Errorf("invalid top-up return URL")
	}

	query := returnURL.Query()
	query.Set("xph_order_no", tradeNo)
	returnURL.RawQuery = query.Encode()
	return returnURL, nil
}
