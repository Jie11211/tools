package filetool

import "github.com/xuri/excelize/v2"

type Excel struct {
	File *excelize.File
}

func NewExcel() *Excel {
	return &Excel{}
}

func (e *Excel) OpenFile(path string) (*excelize.File, error) {
	f, err := excelize.OpenFile(path)
	return f, err
}

func (e *Excel) NewFile() *excelize.File {
	return excelize.NewFile()
}

func (e *Excel) CloseFile(f *excelize.File) error {
	return f.Close()
}

func (e *Excel) SaveFile(f *excelize.File, name ...string) error {
	if len(name) > 0 {
		return f.SaveAs(name[0])
	}
	return f.Save()
}
