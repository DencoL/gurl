package requestcontent

import (
	"gurl/requests"
	"test"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fluentassert/verify"
)

const testFilePath = "../../test/test_files/test_1.hurl"

func TestUpdate_RequestChanged_SendsReadRequestMsg(t *testing.T) {
    model := New()

    _, cmd := model.Update(requests.RequestChanged {
        RequestFilePath: testFilePath,
        IsFolder: false,
    })

    verify.True(test.IsMsgOfType[RequestRead](cmd)).Assert(t, "ReadRequest was not send after changing the request")
}

func TestUpdate_RequestChanged_IsFolder_SetEmptyContent(t *testing.T) {
    model := New()

    newModel, cmd := model.Update(requests.RequestChanged {
        RequestFilePath: testFilePath,
        IsFolder: true,
    })

    verify.False(test.IsMsgOfType[RequestRead](cmd)).Assert(t, "Should not read request content when folder")
    verify.String(newModel.(Model).View()).Empty().Assert(t, "Request content should be empty when folder")
}

func TestUpdate_RequestRead_ChangesContent(t *testing.T) {
    model := New()

    newModel, _ := model.Update(tea.WindowSizeMsg {
        Height: 10,
    })

    newModel, _ = newModel.Update(RequestRead("dummy content"))

    verify.String(newModel.(Model).View()).Contain("dummy content").Assert(t, "Request content is not correct")
}
