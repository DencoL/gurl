package requests

import (
	"testing"
    "test"
    "gurl/data_models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fluentassert/verify"
)

const testPath = "../../test/test_files"

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

    m, _ := model.Update(AllRequestRead([]datamodels.Request {
        {
            Name: "Request 1",
            Method: "GET",
        },
    }))

    newModel := m.(Model)

    verify.Slice(newModel.items.Items()).Should(func(got []list.Item) bool { return len(got) == 1 }).Assert(t)
    verify.String(newModel.items.Items()[0].(datamodels.Request).Name).Equal("Request 1").Assert(t)
    verify.String(newModel.items.Items()[0].(datamodels.Request).Method).Equal("GET").Assert(t)
}

func TestUpdate_AllRequestsReadMsg_SendsChangeRequestMsg(t *testing.T) {
    model := New(testPath)

    _, cmd := model.Update(AllRequestRead([]datamodels.Request {
        {
            Name: "Request 1",
            Method: "GET",
            IsFolder: true,
        },
    }))

    verify.True(test.IsMsgOfType[RequestChanged](cmd)).Assert(t, "RequestChanged was not send after all requests are read")
    verify.String(cmd().(RequestChanged).RequestFilePath).Equal(testPath + "/Request 1.hurl").Assert(t, "Incorrect request path send on change request")
    verify.True(cmd().(RequestChanged).IsFolder).Assert(t, "Incorrect IsFolder send on change request")
}

func TestUpdate_UpAndDownKey_SendsRequestChangedMsg(t *testing.T) {
    model := New(testPath)

    model.items.SetItems([]list.Item {
        datamodels.Request {
            Name: "Request 1",
        },
    })

    keys := []tea.KeyMsg {
        tea.KeyMsg {
            Type: tea.KeyRunes,
            Runes: []rune("j"),
        },
        tea.KeyMsg {
            Type: tea.KeyRunes,
            Runes: []rune("k"),
        },
    }

    for _, key := range keys {
        _, cmd := model.Update(key)


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
        datamodels.Request {
            Name: "Some request",
            IsFolder: false,
        },
    })

    _, cmd := model.Update(tea.KeyMsg(tea.Key {
        Type: tea.KeyEnter,
    }))

    verify.True(test.IsMsgOfType[ExecuteRequest](cmd)).Assert(t, "ExecuteRequest msg was not send after pressing enter")
}

func TestUpdate_EnterKey_IsFolder_ChangesFolder(t *testing.T) {
    model := New(testPath)
    model.items.SetItems([]list.Item {
        datamodels.Request {
            Name: "test_folder",
            IsFolder: true,
        },
    })

    _, cmd := model.Update(tea.KeyMsg(tea.Key {
        Type: tea.KeyEnter,
    }))

    verify.True(test.IsMsgOfType[AllRequestRead](cmd)).Assert(t, "AllRequestRead msg was not send after hitting enter on folder")
    msg, _ := cmd().(AllRequestRead)

    verify.Slice(msg).Should(func(got []datamodels.Request) bool { return len(got) == 2 }).Assert(t)
    verify.String(msg[0].Name).Equal("sub_1").Assert(t)
    verify.String(msg[1].Name).Equal("sub_folder").Assert(t)
}

func TestUpdate_IsFolder_FoldersAreFirst(t *testing.T) {
    model := New(testPath)

    m, _ := model.Update(AllRequestRead([]datamodels.Request {
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
    verify.String(newModel.items.Items()[0].(datamodels.Request).Name).Equal("Folder 2").Assert(t, "Folder should be before requests")
    verify.String(newModel.items.Items()[1].(datamodels.Request).Name).Equal("Folder 1").Assert(t, "Folder should be before requests")
    verify.String(newModel.items.Items()[2].(datamodels.Request).Name).Equal("Request 1").Assert(t)
}

func TestGoBack_RememberRootFolderOnInit(t *testing.T) {
    model := New(testPath)

    verify.Slice(model.folders).Should(func(got []string) bool { return len(got) == 1 }).Assert(t)
    verify.String(model.folders[0]).Equal(testPath + "/").Assert(t)
}

func TestGoBack_EnterOnFolder_RememberFolder(t *testing.T) {
    model := New(testPath)
    model.items.SetItems([]list.Item {
        datamodels.Request {
            Name: "test_folder",
            IsFolder: true,
        },
    })

    newModel, _ := model.Update(tea.KeyMsg(tea.Key {
        Type: tea.KeyEnter,
    }))

    verify.Slice(newModel.(Model).folders).Should(func(got []string) bool { return len(got) == 2 }).Assert(t, "Folder on enter was not remembered")
    verify.String(newModel.(Model).folders[0]).Equal(testPath + "/").Assert(t)
    verify.String(newModel.(Model).folders[1]).Equal(testPath + "/" + "test_folder/").Assert(t)
}

func TestGoBack_GoBackKey_GoesUpOneFolderCorrectly(t *testing.T) {
    model := New(testPath)
    model.folders = append(model.folders, testPath + "/test_folder")
    model.folders = append(model.folders, testPath + "/test_folder/" + "sub_folder")

    newModel, cmd := model.Update(tea.KeyMsg {
        Type: tea.KeyRunes,
        Runes: []rune("-"),
    })

    verify.True(test.IsMsgOfType[AllRequestRead](cmd)).Assert(t, "AllRequestRead msg was not send after go back")
    verify.String(newModel.(Model).folders[0]).Equal(testPath + "/").Assert(t)
    verify.String(newModel.(Model).folders[1]).Equal(testPath + "/test_folder").Assert(t)

    newModel, cmd = model.Update(tea.KeyMsg {
        Type: tea.KeyRunes,
        Runes: []rune("-"),
    })
    
    verify.True(test.IsMsgOfType[AllRequestRead](cmd)).Assert(t, "AllRequestRead msg was not send after go back")
    verify.String(newModel.(Model).folders[0]).Equal(testPath + "/").Assert(t)
}

func TestGoBack_GoBackKey_DoesNotAllowToGoOverRootFolder(t *testing.T) {
    model := New(testPath)

    newModel, cmd := model.Update(tea.KeyMsg {
        Type: tea.KeyRunes,
        Runes: []rune("h"),
    })

    verify.False(test.IsMsgOfType[AllRequestRead](cmd)).Assert(t, "AllRequestRead msg was not send after go back")
    verify.String(newModel.(Model).folders[0]).Equal(testPath + "/").Assert(t)
}
