package validator

import (
	_ "image/png"
)

// TODO: 이미지 파일 포멧에 따라서
// 1. 재인코딩 -> 이미지 파일에 숨기는거 제거
// 2. 유효한 파일인지 확인, 등
// func CheckValidateImage(reader io.Reader) (bytes.Buffer, error) {
// 	img, format, err := image.Decode(reader)
// 	var buf bytes.Buffer

// 	if err != nil {
// 		log.Fatal(err)
// 		return buf, err
// 	}

// 	return buf, err
// }
