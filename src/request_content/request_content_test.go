package requestcontent

import (
	"fmt"
	"gurl/requests"
	"test"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fluentassert/verify"
)

const testFilePath = "../../test/test_files/test_1.hurl"

func TestUpdate_RequestChanged_SendsReadRequestMsg(t *testing.T) {
    model := New()

    _, cmd := model.Update(requests.RequestChanged(testFilePath))

    verify.True(test.IsMsgOfType[RequestRead](cmd)).Assert(t, "ReadRequest was not send after changing the request")
}

func TestUpdate_RequestRead_ChangesContent(t *testing.T) {
    model := New()

    newModel, _ := model.Update(tea.WindowSizeMsg {
        Height: 10,
    })

    newModel, _ = newModel.Update(RequestRead("dummy content"))

    fmt.Println(newModel.View())

    verify.String(newModel.(Model).View()).Contain("dummy content").Assert(t)
}
