package fileProcessor

type tempFile struct {
	path       string
	base64File *string
}

type createdFile struct {
	ID                 string
	Path               string
	Extension          string
	MimeType           string
	PublicFilePath     string
	FileSystemFilePath string
	FileName           string
}

type fileResult struct {
	createdFile createdFile
	error       error
}

type callbackCreateFn = func(fileSystemFilePath, path, mimeType, extension, fileName string) (string, error)
type callbackUpdateFn = func(imageId, fileSystemFilePath, path, mimeType, extension string) error
type callbackDeleteFn = func(imageId, fieldName string) error
