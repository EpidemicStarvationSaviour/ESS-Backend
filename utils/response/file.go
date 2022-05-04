package response

func File(fileName string, fileData []byte) Response {
	return Response{
		Type:     TypeFile,
		File:     fileData,
		FileName: fileName,
	}
}
