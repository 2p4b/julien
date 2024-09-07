package contract

type Driver interface {
	Dump(*map[string]interface{}, string) ([]byte, error)
	Parse([]byte) (*map[string]interface{}, string, error)
}

type Doc interface {
	// Document name
	Name() string

	// Read key from doc
	Get(string, ...interface{}) interface{}

	// Read key from doc
	Has(string) bool

	// Read key from doc
	Body(...string) string
}

type Form interface {
	// Document name
	Name() string

	// Read key from doc
	Has(string) bool

	// Open form doc
	Find(string) (*Doc, error)

	// Open form doc
	Open(string) *Doc

	// Create Form doc
	Compose(string, []byte) (Doc, error)

	// Read form key
	Get(string, ...interface{}) interface{}

	// Set form key value
	Set(string, interface{})

	// Form timestamp
	Body() string
}
