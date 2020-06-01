package dns

import (
	"log"

	"github.com/denverdino/aliyungo/common"
)

type DescribeDomainRecordInfoArgs struct {
	RecordId string
}

type DescribeDomainRecordInfoResponse struct {
	common.Response
	RecordType
}

// DescribeDomainRecordInfo
//
// You can read doc at https://docs.aliyun.com/#/pub/dns/api-reference/record-related&DescribeDomainRecordInfo
func (client *Client) DescribeDomainRecordInfo(args *DescribeDomainRecordInfoArgs) (response *DescribeDomainRecordInfoResponse, err error) {
	action := "DescribeDomainRecordInfo"
	response = &DescribeDomainRecordInfoResponse{}
	err = client.Invoke(action, args, response)
	if err == nil {
		return response, nil
	} else {
		log.Fatalf("%s error, %v", action, err)
		return response, err
	}
}
