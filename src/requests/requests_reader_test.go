package requests

import (
	"gurl/extensions"
	"testing"
    "gurl/data_models"

	"github.com/fluentassert/verify"
)

const testPath = "../../test/test_files"

func TestReadRequestsInfo_EmptyFolder_ReturnsEmpty(t *testing.T) {
    result := ReadRequestsInfo(testPath + "/empty_folder")

    verify.Slice[datamodels.Request](result).Empty().Assert(t)
}

func TestReadRequestsInfo_NonEmptyFolder_ReturnsHurlFiles(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.Slice[datamodels.Request](result).Should(func(got []datamodels.Request) bool { return len(got) == 6 }).Assert(t)

    verify.String(result[0].Name).Equal("empty_folder").Assert(t)
    verify.True(result[0].IsFolder).Assert(t)

    verify.String(result[1].Name).Equal("test_1").Assert(t)
    verify.False(result[1].IsFolder).Assert(t)

    verify.String(result[2].Name).Equal("test_2").Assert(t)
    verify.False(result[2].IsFolder).Assert(t)

    verify.String(result[3].Name).Equal("test_4_no_method").Assert(t)
    verify.False(result[3].IsFolder).Assert(t)

    verify.String(result[4].Name).Equal("test_5_connect").Assert(t)
    verify.False(result[4].IsFolder).Assert(t)

    verify.String(result[5].Name).Equal("test_folder").Assert(t)
    verify.True(result[5].IsFolder).Assert(t)
}

func TestReadRequestsInfo_HttpMethodIsReadFromFile(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.String(result[1].Method).Equal("GET").Assert(t)
    verify.String(result[2].Method).Equal("POST").Assert(t)
    verify.String(result[4].Method).Equal("CONNECT").Assert(t)
}

func TestReadRequestsInfo_HttpMethodIsReadFromFile_UknownMethodMapsToEmptyString(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.String(result[3].Method).Empty().Assert(t)
}

func TestReadRequestsInfo_ReturnsFolders(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.String(result[0].Name).Equal("empty_folder").Assert(t)
    verify.True(result[0].IsFolder).Assert(t)

    verify.String(result[5].Name).Equal("test_folder").Assert(t)
    verify.True(result[5].IsFolder).Assert(t)
}

func TestReadRequestsInfo_EmptyHurlFile_NotAdded(t *testing.T) {
    result := ReadRequestsInfo(testPath)

    verify.Slice[datamodels.Request](result).Should(func(got []datamodels.Request) bool {
        return !extensions.Contains(result, func(request datamodels.Request) bool {
            return request.Name == "test_empty"
        })
    }).Assert(t, "List of requests should not contain empty hurl file")
}
