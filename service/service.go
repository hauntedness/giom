package service

type Service interface {
	Initialized() bool
	Account() Account
	Accounts() <-chan []Account
	Contacts(accountPublicKey string, offset, limit int) <-chan []Contact
	// SetPrivateKey(privateKeyHex string) error
	SetUserPassword(passwd string) <-chan error
	Messages(contactPubKey string, offset, limit int) <-chan []Message
	CreateAccount(privateKeyHex string) <-chan error
	SendMessage(contactPublicKey string, msg string, audioBuf []byte, createdTimestamp string) <-chan error
	SaveContact(contactPublicKey string, identified bool) <-chan error
	AutoCreateAccount() <-chan error
	AccountKeyExists(publicKey string) <-chan bool
	SetAsCurrentAccount(account Account) <-chan error
	Subscribe(topics ...EventTopic) Subscriber
	LastMessage(contactPublicKey string) <-chan Message
	DeleteAccounts([]Account) <-chan error
	DeleteContacts([]Contact) <-chan error
	UnreadMessagesCount(contactPublicKey string) <-chan int64
	MessagesCount(contactPublicKey string) <-chan int64
	ContactsCount(addrPublicKey string) <-chan int64
	AccountsCount() <-chan int64
	MarkPrevMessagesAsRead(contactAddr string) <-chan error
	UserPasswordSet() bool
}

func GetServiceInstance() Service {
	return &service{}
}

var _ Service = (*service)(nil)

type service struct{}

// Account implements Service
func (*service) Account() Account {
	return Account{}
}

// AccountKeyExists implements Service
func (*service) AccountKeyExists(publicKey string) <-chan bool {
	return nil
}

// Accounts implements Service
func (*service) Accounts() <-chan []Account {
	return nil
}

// AccountsCount implements Service
func (*service) AccountsCount() <-chan int64 {
	return nil
}

// AutoCreateAccount implements Service
func (*service) AutoCreateAccount() <-chan error {
	return nil
}

// Contacts implements Service
func (*service) Contacts(accountPublicKey string, offset int, limit int) <-chan []Contact {
	return nil
}

// ContactsCount implements Service
func (*service) ContactsCount(addrPublicKey string) <-chan int64 {
	return nil
}

// CreateAccount implements Service
func (*service) CreateAccount(privateKeyHex string) <-chan error {
	return nil
}

// DeleteAccounts implements Service
func (*service) DeleteAccounts([]Account) <-chan error {
	return nil
}

// DeleteContacts implements Service
func (*service) DeleteContacts([]Contact) <-chan error {
	return nil
}

// Initialized implements Service
func (*service) Initialized() bool {
	return false
}

// LastMessage implements Service
func (*service) LastMessage(contactPublicKey string) <-chan Message {
	return nil
}

// MarkPrevMessagesAsRead implements Service
func (*service) MarkPrevMessagesAsRead(contactAddr string) <-chan error {
	return nil
}

// Messages implements Service
func (*service) Messages(contactPubKey string, offset int, limit int) <-chan []Message {
	return nil
}

// MessagesCount implements Service
func (*service) MessagesCount(contactPublicKey string) <-chan int64 {
	return nil
}

// SaveContact implements Service
func (*service) SaveContact(contactPublicKey string, identified bool) <-chan error {
	return nil
}

// SendMessage implements Service
func (*service) SendMessage(contactPublicKey string, msg string, audioBuf []byte, createdTimestamp string) <-chan error {
	return nil
}

// SetAsCurrentAccount implements Service
func (*service) SetAsCurrentAccount(account Account) <-chan error {
	return nil
}

// SetUserPassword implements Service
func (*service) SetUserPassword(passwd string) <-chan error {
	return nil
}

// Subscribe implements Service
func (*service) Subscribe(topics ...EventTopic) Subscriber {
	return newSubscriber()
}

// UnreadMessagesCount implements Service
func (*service) UnreadMessagesCount(contactPublicKey string) <-chan int64 {
	return nil
}

// UserPasswordSet implements Service
func (*service) UserPasswordSet() bool {
	return false
}
