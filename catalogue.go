package carbon

import (
	"reflect"
	"fmt"
	"errors"
)

type MessageCatalogue struct {
	locale string
	metadata string
	resources map[string]string
	messages map[string]map[string]string
}

func NewMessageCatalogue(locale string, messages map[string]map[string]string) *MessageCatalogue {
	return &MessageCatalogue{
		locale: locale,
		messages: messages,
	}
}

func (mc *MessageCatalogue) GetLocale() string {
	return mc.locale
}

func (mc *MessageCatalogue) All(domain string) map[string]string {
	if domain == "" {
		return mc.messages
	}

	if mc.messages[domain] != nil {
		return mc.messages[domain]
	}

	// TODO Check what to return
	return make(map[string]string)
}

func (mc *MessageCatalogue) Get(id, domain string) string {
	if mc.messages[domain][id] != "" {
		return mc.messages[domain][id]
	}

	//TODO some fallback needed here
	return id
}

func (mc *MessageCatalogue) Add(messages, domain string) {
	if mc.messages[domain] != nil {
		mc.messages[domain] = messages
	} else {
		//TODO do some array replace
		//mc.messages[domain] = array_replace(mc.messages[domain], messages)
	}
}

func (mc *MessageCatalogue) GetResources() map[string]string {
	return mc.resources
}

func (mc *MessageCatalogue) AddResource(resource map[string]string) {
	// TODO
	// $this->resources[$resource->__toString()] = $resource;
	mc.resources[resource] = resource
}

// TODO Check what it returns
func (mc *MessageCatalogue) GetMetadata(key, domain string) string {
	if domain == "" {
		return mc.metadata
	}

	if mc.metadata[domain] != nil {
		if key == nil {
			return mc.metadata[domain]
		}

		if mc.metadata[domain][key] != nil {
			return mc.metadata[domain][key]

		}
	}

	return ""
}

func (mc *MessageCatalogue) addMetadata(values []string) {
	for domain, keys := range values {
		for key, value := range keys {
			mc.SetMetadata(key, value, domain)
		}
	}
}

func (mc *MessageCatalogue) SetMetadata(key, value, domain string) {
	mc.metadata[domain][key] = value
}

func (mc *MessageCatalogue) AddCatalogue(catalogue *MessageCatalogue) error {
	if catalogue.GetLocale() != mc.locale {
		return errors.New(fmt.Sprintf(
			"Cannot add a catalogue for locale %s as the current locale for this catalogue is %s",
			catalogue.GetLocale(),
			mc.locale,
		))
	}

	for i, cat := range catalogue.All("") {
		mc.Add(cat, i)
	}

	for _, resource := range catalogue.GetResources() {
		mc.AddResource(resource)
	}

	// TODO
	// reflect.TypeOf(catalogue) == MetadataAwareInterface
	if reflect.TypeOf(catalogue) == "" {
		metadata := catalogue.GetMetadata("", "")
		mc.addMetadata(metadata)
	}

	return nil
}

// TODO currelny working on this
func  (mc *MessageCatalogue) addFallbackCatalogue(catalogue *MessageCatalogue) {
	// detect circular references
	$c = $catalogue;
	while ($c = $c->getFallbackCatalogue()) {
	if ($c->getLocale() === $this->getLocale()) {
	throw new \LogicException(sprintf('Circular reference detected when adding a fallback catalogue for locale "%s".', $catalogue->getLocale()));
	}
	}

	$c = $this;
	do {
	if ($c->getLocale() === $catalogue->getLocale()) {
	throw new \LogicException(sprintf('Circular reference detected when adding a fallback catalogue for locale "%s".', $catalogue->getLocale()));
	}
	} while ($c = $c->parent);

	$catalogue->parent = $this;
	$this->fallbackCatalogue = $catalogue;

	foreach ($catalogue->getResources() as $resource) {
	$this->addResource($resource);
	}
}
