package err


type Err struct {
	Code string `json:"code"`
	Message string `json:"msg"`
}

func (e *Err) Error() string {
	return e.Message
}

func New(code, message string) error {
	return &Err{Code:code, Message:message}
}

var (
	DBDialectNotSet = New("S9001", "DB dialect not set")
	DBArgsNotSet = New("S9002", "DB args not set")
	DBNotRegistered = New("S9003", "DB not registered")
	DBAlreadyRegistered = New("S9004", "DB already registered")

	NoAPIRegistered = New("S8001", "no API registered")
)