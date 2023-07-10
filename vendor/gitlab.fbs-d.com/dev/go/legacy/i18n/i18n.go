package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

type TranslationsErrors map[string]error //Мы не можем подключить эксепшены, так-как они сами используют i18n

func (err TranslationsErrors) Error() string {
	//Выводит только первую ошибку
	//сделанно по аналогии с MapError
	for k, err := range err {
		return fmt.Sprintf("%s: %s", k, err.Error())
	}
	return ""
}

// map бандлов по категориям переводов
var bundles = &sync.Map{}

// map функций перевода сгруппированных по категориям и языкам, например [errors][en]func
var translations = &sync.Map{}
var defaultLanguage = "en"

// T переводит сообщение
func T(category, message, language string, args ...interface{}) string {
	// Если язык не передан, используем дефолтный
	if language == "" {
		language = defaultLanguage
	}
	// Если есть переводы для этой категории, то пытаемся перевести либо возвращаем перевод на базовом языке
	if t, ok := translations.Load(category); ok {
		if tFunc, ok := t.(*sync.Map).Load(language); ok {
			return tFunc.(bundle.TranslateFunc)(message, args...)
		} else if tFunc, ok := t.(*sync.Map).Load(defaultLanguage); ok {
			return tFunc.(bundle.TranslateFunc)(message, args...)
		}
	}
	return message
}

// SetDefaultLanguage устанавливает дефолтный язык
func SetDefaultLanguage(l string) error {
	langs := language.Parse(l)
	switch l := len(langs); {
	case l == 0:
		return fmt.Errorf("no language found in %q", l)
	case l > 1:
		return fmt.Errorf("multiple languages found in %q: %v; expected one", l, langs)
	}
	defaultLanguage = langs[0].Tag
	return nil
}

// LoadRaw грузит перевод из json переданного в виде строки
func LoadRaw(raw string, category, l string) (err error) {
	// Если язык не передан, используем дефолтный
	if l == "" {
		l = defaultLanguage
	}
	langs := language.Parse(l)
	switch l := len(langs); {
	case l == 0:
		return fmt.Errorf("no language found in %q", l)
	case l > 1:
		return fmt.Errorf("multiple languages found in filename %q: %v; expected one", l, langs)
	}
	lang := langs[0]

	res := make(map[string]string)
	err = json.Unmarshal([]byte(raw), &res)
	if err != nil {
		return
	}
	for key, value := range res {
		// Создаем новый интерфейс Translation из каждой строки файла
		data := map[string]interface{}{
			"id":          key,
			"translation": value,
		}
		var t translation.Translation
		t, err = translation.NewTranslation(data)
		if err != nil {
			return
		}
		err = addTranslation(category, lang, t, false)
		if err != nil {
			return
		}
	}
	return
}

// LoadFromPath загружает все переводы из папки
func LoadFromPath(path string) (err error) {
	var te = make(TranslationsErrors)
	defer func() {
		//В defer мы можем манипулировать возвращаемым значением
		if len(te) > 0 {
			err = te
		}
	}()
	// Считываем все дириктории (для нас это будут языки)
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// Load language directories
	for _, dir := range dirs {
		if !dir.IsDir() {
			return errors.New(fmt.Sprintf("%s is not directory", dir.Name()))
		}

		langs := language.Parse(dir.Name())
		switch l := len(langs); {
		case l == 0:
			return fmt.Errorf("no language found in %q", dir.Name())
		case l > 1:
			return fmt.Errorf("multiple languages found in filename %q: %v; expected one", dir.Name(), langs)
		}
		lang := langs[0]

		dirPath := filepath.Join(path, dir.Name())

		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return err
		}
		// Загружаем все файлы в bundles
		for _, f := range files {
			fp := filepath.Join(dirPath, f.Name())
			ts, err := loadFile(fp)
			if err != nil {
				te[fp] = err
				continue
			}
			name := f.Name()[0:strings.LastIndex(f.Name(), ".json")]

			for _, t := range ts {
				err = addTranslation(name, lang, t, true)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// LoadOneFromPath загружает один файл переводов из папки
func LoadOneFromPath(path, lang, name string) error {
	// парсим язык
	langs := language.Parse(lang)
	if len(langs) != 1 {
		return errors.New("invalid language")
	}

	// загружаем переводы из файла
	filePath := filepath.Join(path, lang, name+".json")
	ts, err := loadFile(filePath)
	if err != nil {
		return err
	}

	for _, t := range ts {
		err = addTranslation(name, langs[0], t, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// addTranslation добавляет перевод в bundle и добавляет функцию перевода
func addTranslation(category string, l *language.Language, t translation.Translation, force bool) (err error) {
	bundleItem, ok := bundles.Load(category)
	if !ok {
		bundleItem = bundle.New()
		bundles.Store(category, bundleItem)
	}
	if !force {
		ts := bundleItem.(*bundle.Bundle).Translations()
		if i, ok := ts[l.Tag]; ok {
			if _, ok := i[t.ID()]; ok {
				return
			}
		}
	}
	bundleItem.(*bundle.Bundle).AddTranslation(l, t)

	var translationCategory *sync.Map
	if tc, ok := translations.Load(category); !ok {
		translationCategory = &sync.Map{}
		translations.Store(category, translationCategory)
	} else {
		translationCategory = tc.(*sync.Map)
	}

	var tag interface{}
	if tag, ok = translationCategory.Load(l.Tag); !ok {
		tag, err = bundleItem.(*bundle.Bundle).Tfunc(l.Tag, defaultLanguage)
		translationCategory.Store(l.Tag, tag)
	}
	return
}

// loadFile загружает json-файл перевода и возвращает массив объектов Translation
func loadFile(filename string) (ts translation.SortableByID, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	res := make(map[string]string)
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return
	}
	ts = make([]translation.Translation, 0)
	for key, value := range res {
		// Создаем новый интерфейс Translation из каждой строки файла
		data := map[string]interface{}{
			"id":          key,
			"translation": value,
		}
		var t translation.Translation
		t, err = translation.NewTranslation(data)
		if err != nil {
			return
		}
		ts = append(ts, t)
	}
	return
}

// GetMessages возвращает массив переводов в виде ключ -> значение
func GetMessages(category, language string) (map[string]string, error) {
	if language == "" {
		language = defaultLanguage
	}
	bundleItem, ok := bundles.Load(category)
	if !ok {
		return nil, errors.New("Unknown category")
	}
	translations := bundleItem.(*bundle.Bundle).Translations()
	if translations[language] == nil {
		return nil, errors.New("Unsupported language")
	}
	messages := make(map[string]string)
	for id, t := range translations[language] {
		messages[id] = t.Template("other").String()
	}
	return messages, nil
}
