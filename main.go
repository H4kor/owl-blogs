package main

import (
	"fmt"
	"owl-blogs/domain/model"
	"reflect"
)

func Persist(entry model.Entry) error {
	t := reflect.TypeOf(entry).Elem().Name()

	fmt.Println(t)
	return nil
}

func main() {
	// repo := infra.NewEntryRepository()
	// repo.RegisterEntryType(&model.ImageEntry{})

	// var img model.Entry = &model.ImageEntry{}
	// img.Create("id", "content", nil, &model.ImageEntryMetaData{ImagePath: "path"})

	// repo.Save(img)

	// img2, err := repo.FindById("id")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(img2)
}
