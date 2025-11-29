package graph

import "time"

const (
	// NodeSettingDefaultMailboxSize is the default size of the node mailbox.
	NodeSettingDefaultMailboxSize = 10
	// NodeSettingDefaultAcceptTimeout is the default timeout for accepting messages in the node mailbox.
	NodeSettingDefaultAcceptTimeout = 5 * time.Second
)

// NodeSettings holds the configuration settings for a graph node.
type NodeSettings struct {
	// MailboxSize is the size of the node mailbox.
	MailboxSize int
	// AcceptTimeout is the timeout for accepting messages in the node mailbox.
	AcceptTimeout time.Duration
}

var defaultNodeSettings = NodeSettings{
	MailboxSize:   NodeSettingDefaultMailboxSize,
	AcceptTimeout: NodeSettingDefaultAcceptTimeout,
}

// FillNodeSettingsWithDefaults fills in any zero-value settings with their default values.
func FillNodeSettingsWithDefaults(s NodeSettings) NodeSettings {
	merged := defaultNodeSettings

	if s.MailboxSize != 0 {
		merged.MailboxSize = s.MailboxSize
	}

	if s.AcceptTimeout != 0 {
		merged.AcceptTimeout = s.AcceptTimeout
	}

	return merged
}
