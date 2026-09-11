package uploadingest

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type fakeVideoStore struct {
	calls  []uuid.UUID
	errors map[uuid.UUID]error
}

func (s *fakeVideoStore) MarkVideoReady(id uuid.UUID) error {
	s.calls = append(s.calls, id)
	return s.errors[id]
}

func TestHandle(t *testing.T) {
	first := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	second := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	record := func(key string) string {
		return fmt.Sprintf(`{"eventSource":"aws:s3","eventName":"ObjectCreated:Put","s3":{"bucket":{"name":"vidtube-videos-1"},"object":{"key":%q}}}`, key)
	}
	firstRecord := record("videos/" + first.String() + "/original.mp4")
	secondRecord := record("videos/" + second.String() + "/original.mp4")
	dbFailure := errors.New("database unavailable")

	tests := []struct {
		name        string
		body        string
		storeErrors map[uuid.UUID]error
		wantAck     bool
		wantCalls   []uuid.UUID
	}{
		{
			name:      "valid upload is acknowledged after updating the video",
			body:      `{"Records":[` + firstRecord + `]}`,
			wantAck:   true,
			wantCalls: []uuid.UUID{first},
		},
		{
			name:        "database failure retains the message for retry",
			body:        `{"Records":[` + firstRecord + `]}`,
			storeErrors: map[uuid.UUID]error{first: dbFailure},
			wantAck:     false,
			wantCalls:   []uuid.UUID{first},
		},
		{
			name:        "record failure retains the message and processing continues",
			body:        `{"Records":[` + firstRecord + `,` + secondRecord + `]}`,
			storeErrors: map[uuid.UUID]error{first: dbFailure},
			wantAck:     false,
			wantCalls:   []uuid.UUID{first, second},
		},
		{
			name:    "S3 test event is acknowledged without updating a video",
			body:    `{"Service":"Amazon S3","Event":"s3:TestEvent","Bucket":"vidtube-videos-1"}`,
			wantAck: true,
		},
		{
			name:    "malformed JSON is acknowledged without updating a video",
			body:    `{"Records":[`,
			wantAck: true,
		},
		{
			name:    "unexpected key layout is acknowledged without updating a video",
			body:    `{"Records":[` + record("other/"+first.String()+"/original.mp4") + `]}`,
			wantAck: true,
		},
		{
			name:    "invalid UUID is acknowledged without updating a video",
			body:    `{"Records":[` + record("videos/not-a-uuid/original.mp4") + `]}`,
			wantAck: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeVideoStore{errors: tt.storeErrors}
			msg := types.Message{MessageId: aws.String("test-message"), Body: aws.String(tt.body)}
			if got := handle(context.Background(), store, msg); got != tt.wantAck {
				t.Errorf("handle() acknowledgment = %v, want %v", got, tt.wantAck)
			}
			if !slices.Equal(store.calls, tt.wantCalls) {
				t.Errorf("MarkVideoReady calls = %v, want %v", store.calls, tt.wantCalls)
			}
		})
	}
}
