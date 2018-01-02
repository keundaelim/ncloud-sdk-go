package server

import (
	"encoding/xml"
	"errors"
	"fmt"

	common "github.com/NaverCloudPlatform/ncloud-sdk-go/common"
	request "github.com/NaverCloudPlatform/ncloud-sdk-go/request"
)

func processTerminateServerInstancesParams(reqParams *RequestTerminateServerInstances) (map[string]string, error) {
	params := make(map[string]string)

	if reqParams == nil || len(reqParams.ServerInstanceNoList) == 0 {
		return params, errors.New("serverInstanceNoList is required")
	}

	if len(reqParams.ServerInstanceNoList) > 0 {
		for k, v := range reqParams.ServerInstanceNoList {
			params[fmt.Sprintf("serverInstanceNoList.%d", k+1)] = v
		}
	}

	return params, nil
}

// TerminateServerInstances terminate server instances
func (s *ServerConn) TerminateServerInstances(reqParams *RequestTerminateServerInstances) (*ServerInstanceList, error) {
	params, err := processTerminateServerInstancesParams(reqParams)
	if err != nil {
		return nil, err
	}

	params["action"] = "terminateServerInstances"

	bytes, resp, err := request.NewRequest(s.accessKey, s.secretKey, "GET", s.apiURL, params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		responseError, err := common.ParseErrorResponse(bytes)
		if err != nil {
			return nil, err
		}

		respError := ServerInstanceList{}
		respError.ReturnCode = responseError.ReturnCode
		respError.ReturnMessage = responseError.ReturnMessage

		return &respError, fmt.Errorf("%s %s", resp.Status, string(bytes))
	}

	var serverInstanceList = ServerInstanceList{}
	if err := xml.Unmarshal([]byte(bytes), &serverInstanceList); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &serverInstanceList, nil
}
