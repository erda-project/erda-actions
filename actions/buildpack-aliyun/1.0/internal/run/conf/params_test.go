package conf

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestHandleParamsBpConfig(t *testing.T) {
	s := "{\"ACL_LINK_URL\":\"https://auth-web.test.cnoocmall.com\",\"DEFAULT_CITYID\":275220,\"DEFAULT_DISTRICTID\":275223,\"DEFAULT_PROVINCEID\":275219,\"DOMAIN_NAME\":\"web-test.mall.cnooc.com.cn\",\"EEVEE_ITEM_ID\":110500400026001,\"EEVEE_SHOP_ID\":4001,\"IM_DOMAIN_NAME\":\"//dev.egoonet.com:60178/zhyIM/webchat.html\",\"LOGISTICS_DOMAIN_NAME\":\"//wuliu.test.cnoocmall.com\",\"OSS_PATH\":\"https://oss.test.cnoocmall.com\",\"PIPELINE_LIMITED_CPU\":1,\"SELLER_DOMAIN_NAME\":\"//seller.test.cnoocmall.com\",\"UC_COOKIE_DOMAIN\":\"test.cnoocmall.com\",\"UPLOAD_ALLOW_HTTP\":false,\"WATERMARK_ENABLE\":true,\"WATERMARK_TEXT\":\"TEST\"}"
	// bp_args
	bpArgs := make(map[string]string)
	// temp bp_args
	tempBpArgs := make(map[string]interface{})
	// bp_args
	err := json.Unmarshal([]byte(s), &tempBpArgs)
	assert.NoError(t, err)
	// map[string]interface{} -> map[string]string
	for k, v := range tempBpArgs {
		switch v.(type) {
		case string:
			bpArgs[k] = v.(string)
		default:
			bpArgs[k] = fmt.Sprintf("%v", v)
		}
	}
	spew.Dump(bpArgs)
}
