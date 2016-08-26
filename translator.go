package carbon

import (
	"errors"
	"regexp"
	"fmt"
	"strings"
)

type Translator struct {
	locale          string
	cacheDir	string
	resources       map[string][]string
	catalogues      map[string]*MessageCatalogue
	fallbackLocales map[string]bool
	loaders         map[string]*Loader
}

// NewTranslator returns a initialized instance of Translator
func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) AddLoader(format string) {
	t.loaders[format] = NewLoader()
}

func (t *Translator) SetLocale(locale string) {
	t.locale = locale
}

func (t *Translator) GetLocale() string {
	return t.locale
}

func (t *Translator) AddResource(locale, format, resource, domain string) error {
	if domain == "" {
		domain = "messages"
	}

	err := t.AssertValidLocale(locale)
	if err != nil {
		return err
	}

	t.resources[locale] = []string{format, resource, domain}

	if t.fallbackLocales[locale] {
		delete(t.catalogues, "locale")
	}

	return nil
}

func (t *Translator) AssertValidLocale(locale string) error {
	matched, err := regexp.MatchString("/^[a-z0-9@_\\.\\-]*$/i", locale)
	if err != nil {
		return errors.New("unable to match locale code : " + err.Error())
	}
	if !matched {
		return errors.New("invalid locale code : " + locale)
	}

	return nil
}

// Trans
func (t *Translator) Trans(id string, parameters []string, domain, locale string) error {
	if domain == "" {
		domain = "messages"
	}

	catalogue, err := t.GetCatalogue(locale)
	if err != nil {
		return err
	}

	//TODO some replacement here
	return (catalogue.Get(id, domain), parameters)
}

// TransChoice
func (t *Translator) TransChoice(id, number string, array parameters = array(), domain = null, locale = null) error {
	if (null === $domain) {
	$domain = 'messages';
	}

	$id = (string) $id;
	$catalogue = $this->getCatalogue($locale);
	$locale = $catalogue->getLocale();
	while (!$catalogue->defines($id, $domain)) {
	if ($cat = $catalogue->getFallbackCatalogue()) {
	$catalogue = $cat;
	$locale = $catalogue->getLocale();
	} else {
	break;
	}
	}

	return strtr($this->selector->choose($catalogue->get($id, $domain), (int) $number, $locale), $parameters);
}

func (t *Translator) GetCatalogue(locale string) (*MessageCatalogue, error) {
	if locale == "" {
		locale = t.GetLocale()
	} else {
		err := t.AssertValidLocale(locale)
		if err != nil {
			return nil, err
		}
	}

	if !t.catalogues[locale] {
		t.LoadCatalogue(locale)
	}

	return t.catalogues[locale], nil
}

func (t *Translator) LoadCatalogue(locale string) {
	if t.cacheDir == "" {
		t.initializeCatalogue(locale)
	} else {
		t.initializeCacheCatalogue(locale)
	}
}

func (t *Translator) DoLoadCatalogue(locale string) error {
	t.catalogues[locale] = NewMessageCatalogue(locale, nil)

	if t.resources[locale] != nil {
		for _, resource := range t.resources[locale] {
			if t.loaders[resource[0]] == nil {
				return errors.New(fmt.Sprintf("The r %s translation loade is not registered : ", resource[0]))
			}
			t.catalogues[locale].AddCatalogue(t.loaders[resource[0]].Load(resource[1], locale, resource[2]))
		}
	}
}

func (t *Translator) computeFallbackLocales(locale string) []string {
	locales := []string{}
	for fallback := range t.fallbackLocales {
		if fallback == locale {
			continue
		}

		locales = append(locales, fallback)
	}

	if strings.Contains(locale, "_") != false {
		locales = append(locale[0 : -len(strings.Contains(locale, "_"))], locales)
	}

	return locales
}

func (t *Translator) loadFallbackCatalogues(locale string) error {
	current := t.catalogues[locale]

	for _, fallback := range t.computeFallbackLocales(locale) {
		if t.catalogues[fallback] == nil {
			err := t.doLoadCatalogue(fallback)
			if err != nil {
				return err
			}
		}

		fallbackCatalogue := NewMessageCatalogue(fallback, t.catalogues[fallback].All())
		for _, resource := range t.catalogues[fallback].GetResources() {
			fallbackCatalogue.AddResource(resource)
		}
		current.AddFallbackCatalogue(fallbackCatalogue)
		current = fallbackCatalogue
	}
}

func (t *Translator) doLoadCatalogue(locale string) error {
	t.catalogues[locale] = NewMessageCatalogue(locale, nil)

	if t.resources[locale] != nil {
		for _, resource := range t.resources[locale] {
			if t.loaders[resource[0]] == nil {
				return  errors.New(fmt.Sprintf("The %s translation loader is not registered", resource[0]))
			}
			t.catalogues[locale].AddCatalogue(t.loaders[resource[0]].Load(resource[1], locale, resource[2]))
		}
	}
}

func (t *Translator) initializeCatalogue(locale string) error {
	err := t.AssertValidLocale(locale)
	if err != nil {
		return err
	}

	err = t.DoLoadCatalogue(locale)
	if err != nil {
		if !t.computeFallbackLocales(locale) {
			return err
		}
	}

	t.loadFallbackCatalogues(locale)

	return nil
}

func (t *Translator) initializeCacheCatalogue(locale string) {
	if (isset($this->catalogues[$locale])) {
	/* Catalogue already initialized. */
	return;
	}

	$this->assertValidLocale($locale);
	$self = $this; // required for PHP 5.3 where "$this" cannot be use()d in anonymous functions. Change in Symfony 3.0.
	$cache = $this->getConfigCacheFactory()->cache($this->getCatalogueCachePath($locale),
	function (ConfigCacheInterface $cache) use ($self, $locale) {
	$self->dumpCatalogue($locale, $cache);
	}
	);

	if (isset($this->catalogues[$locale])) {
	/* Catalogue has been initialized as it was written out to cache. */
	return;
	}

	/* Read catalogue from cache. */
	$this->catalogues[$locale] = include $cache->getPath();
}
