package requests

import (
	"testing"
    "test"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fluentassert/verify"
)

func TestUpdate_WindowSizeMsg_SendsAllRequestReadMsg(t *testing.T) {
    model := New(testPath)

    _, cmd := model.Update(tea.WindowSizeMsg {
        Width: 100,
        Height: 100,
    })

    verify.True(test.IsMsgOfType[AllRequestRead](cmd)).Assert(t, "AllRequestRead msg was not send after WindowSizeMsg")
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
            IsFolder: true,
        },
    }))

    verify.True(test.IsMsgOfType[RequestChanged](cmd)).Assert(t, "RequestChanged was not send after all requests are read")
    verify.String(cmd().(RequestChanged).RequestFilePath).Equal(testPath + "Request 1.hurl").Assert(t, "Incorrect request path send on change request")
    verify.True(cmd().(RequestChanged).IsFolder).Assert(t, "Incorrect IsFolder send on change request")
}

func TestUpdate_UpAndDownKey_SendsRequestChangedMsg(t *testing.T) {
    model := New(testPath)

    model.items.SetItems([]list.Item {
        Request {
            Name: "Request 1",
        },
    })

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
            return test.IsMsgOfType[RequestChanged](c)
        }).Assert(t, "RequestChanged msg was not send after moving up/down in the list")
    }

}

func TestUpdate_EnterKey_SendsExecuteRequestMsg(t *testing.T) {
    model := New(testPath + "test_1.hurl")
    model.items.SetItems([]list.Item {
        Request {
            Name: "Some request",
            IsFolder: false,
        },
    })

    _, cmd := model.Update(tea.KeyMsg(tea.Key {
        Type: tea.KeyEnter,
    }))

    verify.True(test.IsMsgOfType[ExecuteRequest](cmd)).Assert(t, "ExecuteRequest msg was not send after pressing enter")
}

func TestUpdate_EnterKey_IsFolder_DoNothing(t *testing.T) {
    model := New(testPath)
    model.items.SetItems([]list.Item {
        Request {
            Name: "folder",
            IsFolder: true,
        },
    })

    _, cmd := model.Update(tea.KeyMsg(tea.Key {
        Type: tea.KeyEnter,
    }))

    verify.Any[tea.Cmd](cmd).Should(func(got tea.Cmd) bool { return got == nil }).Assert(t, "Enter should do nothing on folder")
}

func TestUpdate_IsFolder_FoldersAreFirst(t *testing.T) {
    model := New(testPath)

    m, _ := model.Update(AllRequestRead([]Request {
        {
            Name: "Folder 2",
            IsFolder: true,
        },
        {
            Name: "Request 1",
        },
        {
            Name: "Folder 1",
            IsFolder: true,
        },
    }))

    newModel := m.(Model)

    verify.Slice(newModel.items.Items()).Should(func(got []list.Item) bool { return len(got) == 3 }).Assert(t)
    verify.String(newModel.items.Items()[0].(Request).Name).Equal("Folder 2").Assert(t)
    verify.String(newModel.items.Items()[1].(Request).Name).Equal("Folder 1").Assert(t)
    verify.String(newModel.items.Items()[2].(Request).Name).Equal("Request 1").Assert(t)
}
