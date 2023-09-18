package requestresponse

import (
	"gurl/requests"
	"test"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fluentassert/verify"
)

const testFilePath = "../../test/test_files/test_1.hurl"

func TestUpdate_ExecuteRequest_SendsRequestExecutedMsg(t *testing.T) {
    model := New()

    _, cmd := model.Update(requests.ExecuteRequest(testFilePath))

    verify.True(test.IsMsgOfType[RequestExecuted](cmd)).Assert(t, "RequestExecuted was not sent after receiving execute request")
}

func TestUpdate_RequestExecuted_ContentIsSet(t *testing.T) {
    model := New()

    newModel, _ := model.Update(tea.WindowSizeMsg {
        Height: 10,
    })

    newModel, _ = newModel.Update(RequestExecuted("dummy response"))

    verify.String(newModel.(Model).View()).Contain("dummy response").Assert(t, "Request response is not correct")
}

func TestUpdate_RequestExecuted_EmptyResponse_SetsEmptyResponseString(t *testing.T) {
    model := New()

    newModel, _ := model.Update(tea.WindowSizeMsg {
        Height: 10,
    })

    newModel, _ = newModel.Update(RequestExecuted(""))

    verify.String(newModel.(Model).View()).Contain("<EMPTY RESPONSE>").Assert(t, "Request response is not empty")
}
