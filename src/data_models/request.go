package datamodels

type Request struct {
	Name     string
	Method   string
	IsFolder bool
}

func (self Request) FilterValue() string {
	return self.Name
}

func (self Request) Title() string {
	if self.IsFolder {
		return "/" + self.Name
	}

	return self.Name
}

func (self Request) Description() string {
	if self.IsFolder {
		return "FOLDER"
	}

	return self.Method
}
