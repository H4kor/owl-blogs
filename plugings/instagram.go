package plugings

import (
	"bytes"
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	entrytypes "owl-blogs/entry_types"
	"owl-blogs/render"

	"github.com/Davincible/goinsta/v3"
)

type Instagram struct {
	configRepo repository.ConfigRepository
	binService *app.BinaryService
}

type InstagramConfig struct {
	User     string
	Password string
}

// Form implements app.AppConfig.
func (cfg *InstagramConfig) Form(binSvc model.BinaryStorageInterface) string {
	f, _ := render.RenderTemplateToString("forms/InstagramConfig", cfg)
	return f
}

// ParseFormData implements app.AppConfig.
func (*InstagramConfig) ParseFormData(data model.HttpFormData, binSvc model.BinaryStorageInterface) (app.AppConfig, error) {
	return &InstagramConfig{
		User:     data.FormValue("User"),
		Password: data.FormValue("Password"),
	}, nil
}

func RegisterInstagram(
	configRepo repository.ConfigRepository,
	configRegister *app.ConfigRegister,
	binService *app.BinaryService,
	bus *app.EventBus,
) *Instagram {
	configRegister.Register("instagram", &InstagramConfig{})
	insta := &Instagram{
		configRepo: configRepo,
		binService: binService,
	}

	bus.Subscribe(insta)

	return insta
}

// NotifyEntryCreated implements app.EntryCreationSubscriber.
func (i *Instagram) NotifyEntryCreated(entry model.Entry) {

	image, ok := entry.(*entrytypes.Image)
	if !ok {
		println("not an image")
		return
	}

	config := &InstagramConfig{}
	err := i.configRepo.Get("instagram", config)
	if err != nil {
		println("no instagram config")
		return
	}

	client := goinsta.New(config.User, config.Password)

	err = client.Login()
	if err != nil {
		println("login failed")
		return
	}

	meta := image.MetaData().(*entrytypes.ImageMetaData)
	bin, err := i.binService.FindById(meta.ImageId)
	if err != nil {
		println("image data not found")
		return
	}

	_, err = client.Upload(
		&goinsta.UploadOptions{
			File:    bytes.NewReader(bin.Data),
			Caption: image.Title(),
		},
	)
	if err != nil {
		println("upload failed")
		return
	}

}
