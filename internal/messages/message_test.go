package messages_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/asankov/secure-messenger/internal/messages"
	"github.com/stretchr/testify/require"
)

const (
	senderID   = "ABC123"
	receiverID = "XYZ987"
	payload    = "Hello there"
)

var (
	timestamp = time.Now().UTC().Unix()
)

func TestNewMessage(t *testing.T) {
	now := time.Now().UTC().Unix()
	msg, err := messages.NewMessage(senderID, receiverID, payload)

	require.NoError(t, err)
	require.Equal(t, senderID, msg.SenderID)
	require.Equal(t, receiverID, msg.ReceiverID)
	require.Equal(t, payload, msg.Payload)
	// since there will always be some difference in the timestamp
	// compare them and allow them to be different by no more than 2 seconds
	timeDiff := now - msg.Timestamp
	if timeDiff > 2 || timeDiff < -2 {
		t.Errorf("difference between time.Now in the test and the timestamp of the message is more than 2 seconds")
	}
}

func TestNewMessageValidation(t *testing.T) {
	t.Run("EmptySenderID", func(t *testing.T) {
		_, err := messages.NewMessage("", receiverID, payload)
		require.ErrorIs(t, err, messages.ErrEmptySenderID)
	})
	t.Run("EmptyReceiverID", func(t *testing.T) {
		_, err := messages.NewMessage(senderID, "", payload)
		require.ErrorIs(t, err, messages.ErrEmptyReceiverID)
	})
	t.Run("EmptyPayload", func(t *testing.T) {
		_, err := messages.NewMessage(senderID, receiverID, "")
		require.ErrorIs(t, err, messages.ErrEmptyPayload)
	})
}

func TestParseFromJSON(t *testing.T) {
	expectedTimestamp := time.Now().UTC().Unix()
	json := fmt.Sprintf(`{"senderId":"%s", "receiverId":"%s", "payload":"%s", "timestamp":%d}`, senderID, receiverID, payload, expectedTimestamp)

	msg, err := messages.ParseFromJSON(json)
	require.NoError(t, err)

	require.Equal(t, senderID, msg.SenderID)
	require.Equal(t, receiverID, msg.ReceiverID)
	require.Equal(t, payload, msg.Payload)
	require.Equal(t, expectedTimestamp, msg.Timestamp)
}

func TestParseFromJSONValidation(t *testing.T) {
	testCases := []struct {
		name          string
		json          string
		expectedError error
	}{
		{
			name:          "no senderId",
			json:          fmt.Sprintf(`{"receiverId":"%s", "payload":"%s", "timestamp":%d}`, receiverID, payload, timestamp),
			expectedError: messages.ErrEmptySenderID,
		},
		{
			name:          "no receiverId",
			json:          fmt.Sprintf(`{"senderId":"%s", "payload":"%s", "timestamp":%d}`, senderID, payload, timestamp),
			expectedError: messages.ErrEmptyReceiverID,
		},
		{
			name:          "no payload",
			json:          fmt.Sprintf(`{"senderId":"%s", "receiverId":"%s","timestamp":%d}`, senderID, receiverID, timestamp),
			expectedError: messages.ErrEmptyPayload,
		},
		{
			name:          "no timestamp",
			json:          fmt.Sprintf(`{"senderId":"%s", "receiverId":"%s", "payload":"%s"}`, senderID, receiverID, payload),
			expectedError: messages.ErrEmptyTimestamp,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := messages.ParseFromJSON(testCase.json)
			require.ErrorIs(t, err, testCase.expectedError)
		})
	}

	// this test does not test much, but it makes the test cover 100% of the code
	t.Run("empty message", func(t *testing.T) {
		_, err := messages.ParseFromJSON("")
		require.Error(t, err)
	})
}
