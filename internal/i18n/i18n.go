package i18n

import (
	"embed"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/Xuanwo/go-locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localeFS embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
)

var SupportedLocales = []language.Tag{
	language.French,
	language.English,
}

func Init() error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, locale := range SupportedLocales {
		_, err := bundle.LoadMessageFileFS(localeFS, "locales/"+locale.String()+".toml")
		if err != nil {
			return fmt.Errorf("failed to load %s locale: %w", locale.String(), err)
		}
	}

	tag, err := locale.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect locale: %w", err)
	}

	localizer = i18n.NewLocalizer(bundle, tag.String())

	return nil
}

func T(messageID string, templateData map[string]any) string {
	if localizer == nil {
		return messageID
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		return messageID
	}

	return msg
}
