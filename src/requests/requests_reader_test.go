package requests

import (
	"testing"
	"github.com/fluentassert/verify"
)

const testPath = "../../test/test_files"

func TestReadRequestsInfo_EmptyFolder_ReturnsEmpty(t *testing.T) {
    result := ReadRequestsInfo(testPath + "/empty_folder")

    verify.Slice[RequestInfo](result).Empty().Assert(t)
}

func TestReadRequestsInfo_NonEmptyFolder_ReturnsHurlFiles(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.Slice[RequestInfo](result).Should(func(got []RequestInfo) bool { return len(got) == 3 }).Assert(t)
    verify.String(result[0].name).Equal("test_1.hurl").Assert(t)
    verify.String(result[1].name).Equal("test_2.hurl").Assert(t)
}

func TestReadRequestsInfo_HttpMethodIsReadFromFile(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.String(result[0].method).Equal("GET").Assert(t)
    verify.String(result[1].method).Equal("POST").Assert(t)
}

func TestReadRequestsInfo_HttpMethodIsReadFromFile_UknownMethodMapsToEmptyString(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.String(result[2].method).Empty().Assert(t)
}
