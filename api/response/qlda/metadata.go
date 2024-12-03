package qlda

type FinalMetadata struct {
	Metadata Metadata `json:"metadata"`
	Pnum     int      `json:"pnum"`
	Page     []Page   `json:"page"`
	Message  string   `json:"message"`
}

type Metadata struct {
	Languages string   `json:"languages"`
	Pages     int      `json:"pages"`
	FileType  string   `json:"fileType"`
	PdfToc    []PdfToc `json:"pdf_toc"`
}
type PdfToc struct {
	Title string `json:"title"`
	Level int    `json:"level"`
	Page  int    `json:"page"`
}

type Page struct {
	Page   int           `json:"page"`
	Text   string        `json:"text"`
	Images []interface{} `json:"images"`
	Tables []interface{} `json:"tables"`
}

func MockFinalMetadata() FinalMetadata {
	return FinalMetadata{
		Metadata: Metadata{
			Languages: "English",
			Pages:     10,
			FileType:  "PDF",
			PdfToc: []PdfToc{
				{Title: "Introduction", Level: 1, Page: 1},
				{Title: "Chapter 1", Level: 1, Page: 2},
			},
		},
		Pnum: 1,
		Page: []Page{
			{Page: 1, Text: "Sample text", Images: []interface{}{"image1.jpg"}, Tables: []interface{}{"table1"}},
		},
		Message: "Mock data generated",
	}
}
