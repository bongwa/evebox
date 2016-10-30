package agent

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSenderQueue(t *testing.T) {
	sender := NewEventSender()
	assert.Equal(t, cap(sender.queue), 1000)
	assert.Equal(t, len(sender.queue), 0)

	sender.Submit("some event")
	assert.Equal(t, len(sender.queue), 1)

	sender.Reset()
	assert.Equal(t, cap(sender.queue), 1000)
	assert.Equal(t, len(sender.queue), 0)
}
