//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"

	"github.com/google/uuid"
)

// FromContext allows for the retrieval of the specified key's value from the supplied Context.
// If the value is not found, an empty string is returned.
func FromContext(ctx context.Context, key string) string {
	hdr, ok := ctx.Value(key).(string)
	if !ok {
		hdr = ""
	}
	return hdr
}

// correlatedId gets Correlation ID from supplied context. If no Correlation ID header is
// present in the supplied context, one will be created along with a value.
func correlatedId(ctx context.Context) string {
	correlation := FromContext(ctx, common.CorrelationHeader)
	if len(correlation) == 0 {
		correlation = uuid.New().String()
	}
	return correlation
}

// Helper method to get the body from the response after making the request
func getBody(resp *http.Response) ([]byte, errors.EdgeX) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, errors.NewCommonEdgeX(errors.KindIOError, "读取响应失败", err)
	}
	return body, nil
}

// Helper method to make the request and return the response
func makeRequest(req *http.Request) (*http.Response, errors.EdgeX) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "发起请求失败", err)
	}
	if resp == nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "响应不能为空", nil)
	}
	return resp, nil
}

func createRequest(ctx context.Context, httpMethod string, baseUrl string, requestPath string, requestParams url.Values) (*http.Request, errors.EdgeX) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "解析URL失败", err)
	}
	u.Path = requestPath
	if requestParams != nil {
		u.RawQuery = requestParams.Encode()
	}
	req, err := http.NewRequest(httpMethod, u.String(), nil)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "创建HTTP请求失败", err)
	}
	req.Header.Set(common.CorrelationHeader, correlatedId(ctx))
	return req, nil
}

func createRequestWithRawDataAndParams(ctx context.Context, httpMethod string, baseUrl string, requestPath string, requestParams url.Values, data interface{}) (*http.Request, errors.EdgeX) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "解析URL失败", err)
	}
	u.Path = requestPath
	if requestParams != nil {
		u.RawQuery = requestParams.Encode()
	}
	jsonEncodedData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "序列化为JSON失败", err)
	}

	content := FromContext(ctx, common.ContentType)
	if content == "" {
		content = common.ContentTypeJSON
	}

	req, err := http.NewRequest(httpMethod, u.String(), bytes.NewReader(jsonEncodedData))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "创建HTTP请求失败", err)
	}
	req.Header.Set(common.ContentType, content)
	req.Header.Set(common.CorrelationHeader, correlatedId(ctx))
	return req, nil
}

func createRequestWithRawData(ctx context.Context, httpMethod string, url string, data interface{}) (*http.Request, errors.EdgeX) {
	jsonEncodedData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "序列化为JSON失败", err)
	}

	content := FromContext(ctx, common.ContentType)
	if content == "" {
		content = common.ContentTypeJSON
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(jsonEncodedData))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "创建HTTP请求失败", err)
	}
	req.Header.Set(common.ContentType, content)
	req.Header.Set(common.CorrelationHeader, correlatedId(ctx))
	return req, nil
}

func createRequestWithEncodedData(ctx context.Context, httpMethod string, url string, data []byte, encoding string) (*http.Request, errors.EdgeX) {
	content := encoding
	if content == "" {
		content = FromContext(ctx, common.ContentType)
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "创建HTTP请求失败", err)
	}
	req.Header.Set(common.ContentType, content)
	req.Header.Set(common.CorrelationHeader, correlatedId(ctx))
	return req, nil
}

// createRequestFromFilePath creates multipart/form-data request with the specified file
func createRequestFromFilePath(ctx context.Context, httpMethod string, url string, filePath string) (*http.Request, errors.EdgeX) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindIOError, fmt.Sprintf("读取文件失败： %s", filePath), err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFileWriter, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "创建表单数据失败", err)
	}
	_, err = io.Copy(formFileWriter, bytes.NewReader(fileContents))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindIOError, "从表单复制数据失败", err)
	}
	writer.Close()

	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "创建HTTP请求失败", err)
	}
	req.Header.Set(common.ContentType, writer.FormDataContentType())
	req.Header.Set(common.CorrelationHeader, correlatedId(ctx))
	return req, nil
}

// sendRequest will make a request with raw data to the specified URL.
// It returns the body as a byte array if successful and an error otherwise.
func sendRequest(ctx context.Context, req *http.Request) ([]byte, errors.EdgeX) {
	resp, err := makeRequest(req)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := getBody(resp)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}

	if resp.StatusCode <= http.StatusMultiStatus {
		return bodyBytes, nil
	}

	var errMessage string
	baseResp := dtoCommon.BaseResponse{}
	if err := json.Unmarshal(bodyBytes, &baseResp); err == nil {
		errMessage = baseResp.Message
	} else {
		errMessage = fmt.Sprintf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	errKind := errors.KindMapping(resp.StatusCode)
	return nil, errors.NewCommonEdgeX(errKind, errMessage, nil)
}
