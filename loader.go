package carbon

// Loader is the structure for loading
type Loader struct {
}

// NewLoader returns a initialized instance of Loader
func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(resource, locale, domain string) *MessageCatalogue {
	//TODO send proper data
	l.flatten(resource, "", "")
	catalogue := NewMessageCatalogue(locale, nil)
	catalogue.Add(resource, domain)

	return catalogue
}

func (l *Loader) flatten(messages map[string]string, subnode, path string) {
	if subnode == nil {
		subnode = messages
	}
	for key, value := range subnode {
		if value != nil {
			var nodePath string
			if path != "" {
				nodePath = path + "." + key
			} else {
				nodePath = key
			}

			l.flatten(messages, value, nodePath)
			if path == "" {
				delete(messages, key)
			}
		} else if path != "" {
			messages[path+"."+key] = value
		}
	}
}
