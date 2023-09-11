package requests

type Request struct {
    Name string
    Method string
}

func (self Request) FilterValue() string {
    return self.Name
}

func (self Request) Title() string {
    return self.Name
}

func (self Request) Description() string {
    return self.Method
}
