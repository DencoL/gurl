package requests

import (
	"testing"
	"github.com/fluentassert/verify"
)

func TestReadRequestsInfo_EmptyFolder_ReturnsEmpty(t *testing.T) {
    result := ReadRequestsInfo("../test/test_files/empty_folder")

    verify.Slice[RequestInfo](result).Empty().Assert(t)
}

func TestReadRequestsInfo_NonEmptyFolder_ReturnsHurlFiles(t *testing.T) {
    result := ReadRequestsInfo("../test/test_files/")

    verify.Slice[RequestInfo](result).Should(func(got []RequestInfo) bool { return len(got) == 2 }).Assert(t)
    verify.String(result[0].name).Equal("test_1.hurl")
    verify.String(result[1].name).Equal("test_2.hurl")
}
