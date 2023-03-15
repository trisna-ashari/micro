package object

// Option return Metadata with Option.
type Option func(metadata *Metadata)

// WithID is a function uses to set Metadata.ID.
func WithID(id string) Option {
	return func(m *Metadata) {
		m.ID = id
	}
}

// WithToken is a function uses to set Metadata.Token.
func WithToken(id string) Option {
	return func(m *Metadata) {
		m.Token = id
	}
}

// IncludeSlug is a function uses to include Metadata.Slug to the path.
func IncludeSlug() Option {
	return func(m *Metadata) {
		m.IncludeSlug = true
	}
}

// IncludeDate is a function uses to include Metadata.Date to the path.
func IncludeDate() Option {
	return func(m *Metadata) {
		m.IncludeDate = true
	}
}

// WithName is a function uses to set Metadata.Name.
func WithName(name string) Option {
	return func(m *Metadata) {
		m.Name = name
	}
}

// WithPrefixOnFileName is a function uses to set Metadata.NamePrefix.
func WithPrefixOnFileName(namePrefix string) Option {
	return func(m *Metadata) {
		m.NamePrefix = namePrefix
	}
}

// WithSuffixOnFileName is a function uses to set Metadata.NameSuffix.
func WithSuffixOnFileName(nameSuffix string) Option {
	return func(m *Metadata) {
		m.NameSuffix = nameSuffix
	}
}

// WithPassword is a function uses to set Metadata.Password.
func WithPassword(password string) Option {
	return func(m *Metadata) {
		m.Password = password
	}
}

// WithPDFOverwriteMode is a function uses to set Metadata.PDFOverwrite.
func WithPDFOverwriteMode(overwriteMode PDFOverwriteMode) Option {
	return func(m *Metadata) {
		m.PDFOverwrite = overwriteMode
	}
}

// WithSource is a function to set Metadata.Source.
func WithSource(source string) Option {
	return func(m *Metadata) {
		m.Source = source
	}
}
