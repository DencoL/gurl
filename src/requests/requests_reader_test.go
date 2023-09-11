package requests

import (
	"testing"
	"github.com/fluentassert/verify"
)

const testPath = "../../test/test_files"

func TestReadRequestsInfo_EmptyFolder_ReturnsEmpty(t *testing.T) {
    result := ReadRequestsInfo(testPath + "/empty_folder")

    verify.Slice[Request](result).Empty().Assert(t)
}

func TestReadRequestsInfo_NonEmptyFolder_ReturnsHurlFiles(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.Slice[Request](result).Should(func(got []Request) bool { return len(got) == 4 }).Assert(t)
    verify.String(result[0].Name).Equal("test_1").Assert(t)
    verify.String(result[1].Name).Equal("test_2").Assert(t)
    verify.String(result[2].Name).Equal("test_4_no_method").Assert(t)
    verify.String(result[3].Name).Equal("test_5_connect").Assert(t)
}

func TestReadRequestsInfo_HttpMethodIsReadFromFile(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.String(result[0].Method).Equal("GET").Assert(t)
    verify.String(result[1].Method).Equal("POST").Assert(t)
    verify.String(result[3].Method).Equal("CONNECT").Assert(t)
}

func TestReadRequestsInfo_HttpMethodIsReadFromFile_UknownMethodMapsToEmptyString(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.String(result[2].Method).Empty().Assert(t)
}
