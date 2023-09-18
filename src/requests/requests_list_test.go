package requests

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fluentassert/verify"
)

func isMsgOfType[TCommand tea.Msg](cmd tea.Cmd) bool {
    if cmd == nil {
        return false
    }

    _, ok := cmd().(TCommand)

    return ok
}

func TestUpdate_WindowSizeMsg_SendsAllRequestReadMsg(t *testing.T) {
    model := New(testPath)

    _, cmd := model.Update(tea.WindowSizeMsg {
        Width: 100,
        Height: 100,
    })

    verify.True(isMsgOfType[AllRequestRead](cmd)).Assert(t, "AllRequestRead msg was not send after WindowSizeMsg")
}

func TestUpdate_AllRequestsReadMsg_SetsReceivedRequests(t *testing.T) {
    model := New(testPath)

    m, _ := model.Update(AllRequestRead([]Request {
        {
            Name: "Request 1",
            Method: "GET",
        },
    }))

    newModel := m.(Model)

    verify.Slice(newModel.items.Items()).Should(func(got []list.Item) bool { return len(got) == 1 }).Assert(t)
    verify.String(newModel.items.Items()[0].(Request).Name).Equal("Request 1").Assert(t)
    verify.String(newModel.items.Items()[0].(Request).Method).Equal("GET").Assert(t)
}

func TestUpdate_AllRequestsReadMsg_SendsChangeRequestMsg(t *testing.T) {
    model := New(testPath)

    _, cmd := model.Update(AllRequestRead([]Request {
        {
            Name: "Request 1",
            Method: "GET",
        },
    }))

    verify.True(isMsgOfType[RequestChanged](cmd)).Assert(t, "RequestChanged was not send after all requests are read")
}

func TestUpdate_UpAndDownKey_SendsRequestChangedMsg(t *testing.T) {
    model := New(testPath)


    keys := []tea.KeyType {
        tea.KeyDown,
        tea.KeyUp,
    }

    for _, key := range keys {
        _, cmd := model.Update(tea.KeyMsg(tea.Key {
            Type: key,
        }))


        var cmds tea.BatchMsg
        verify.NotPanics(func() {
            cmds = cmd().(tea.BatchMsg)
        }).Assert(t, "Cmd is not of type BatchMsg")

        verify.Slice[tea.Cmd](cmds).Any(func(c tea.Cmd) bool {
            return isMsgOfType[RequestChanged](c)
        }).Assert(t, "RequestChanged msg was not send after moving up/down in the list")
    }

}

func TestUpdate_EnterKey_SendsExecuteRequestMsg(t *testing.T) {
    model := New(testPath)

    _, cmd := model.Update(tea.KeyMsg(tea.Key {
        Type: tea.KeyEnter,
    }))

    verify.True(isMsgOfType[ExecuteRequest](cmd)).Assert(t, "ExecuteRequest msg was not send after pressing enter")
}
