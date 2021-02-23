package bot

import (
	"testing"
)

func TestConfig_ShouldReplyUnsure(t *testing.T) {
	type fields struct {
		Conversation Conversation
	}
	type args struct {
		isExistingConversation bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should ReplyUnsure to a new conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyUnsure: true,
					},
					Existing: ConversationConfig{
						ReplyUnsure: false,
					},
				},
			},
			args: args{isExistingConversation: false},
			want: true,
		},
		{
			name: "should ReplyUnsure to an existing conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyUnsure: false,
					},
					Existing: ConversationConfig{
						ReplyUnsure: true,
					},
				},
			},
			args: args{isExistingConversation: true},
			want: true,
		},
		{
			name: "should not ReplyUnsure to a new conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyUnsure: false,
					},
					Existing: ConversationConfig{
						ReplyUnsure: false,
					},
				},
			},
			args: args{isExistingConversation: false},
			want: false,
		},
		{
			name: "should not ReplyUnsure to an existing conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyUnsure: false,
					},
					Existing: ConversationConfig{
						ReplyUnsure: false,
					},
				},
			},
			args: args{isExistingConversation: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Conversation: tt.fields.Conversation,
			}
			if got := c.ShouldReplyUnsure(tt.args.isExistingConversation); got != tt.want {
				t.Errorf("ShouldReplyUnsure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_ShouldReplyUnknown(t *testing.T) {
	type fields struct {
		Conversation Conversation
	}
	type args struct {
		isExistingConversation bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should ReplyUnknown to a new conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyUnknown: true,
					},
					Existing: ConversationConfig{
						ReplyUnknown: false,
					},
				},
			},
			args: args{isExistingConversation: false},
			want: true,
		},
		{
			name: "should ReplyUnknown to an existing conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyUnknown: false,
					},
					Existing: ConversationConfig{
						ReplyUnknown: true,
					},
				},
			},
			args: args{isExistingConversation: true},
			want: true,
		},
		{
			name: "should not ReplyUnknown to a new conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyUnknown: false,
					},
					Existing: ConversationConfig{
						ReplyUnknown: false,
					},
				},
			},
			args: args{isExistingConversation: false},
			want: false,
		},
		{
			name: "should not ReplyUnknown to an existing conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyUnknown: false,
					},
					Existing: ConversationConfig{
						ReplyUnknown: false,
					},
				},
			},
			args: args{isExistingConversation: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Conversation: tt.fields.Conversation,
			}
			if got := c.ShouldReplyUnknown(tt.args.isExistingConversation); got != tt.want {
				t.Errorf("ShouldReplyUnknown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_ShouldReplyError(t *testing.T) {
	type fields struct {
		Conversation Conversation
	}
	type args struct {
		isExistingConversation bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should ReplyError to a new conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyError: true,
					},
					Existing: ConversationConfig{
						ReplyError: false,
					},
				},
			},
			args: args{isExistingConversation: false},
			want: true,
		},
		{
			name: "should ReplyError to an existing conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyError: false,
					},
					Existing: ConversationConfig{
						ReplyError: true,
					},
				},
			},
			args: args{isExistingConversation: true},
			want: true,
		},
		{
			name: "should not ReplyError to a new conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyError: false,
					},
					Existing: ConversationConfig{
						ReplyError: false,
					},
				},
			},
			args: args{isExistingConversation: false},
			want: false,
		},
		{
			name: "should not ReplyError to an existing conversation",
			fields: fields{
				Conversation: Conversation{
					New: ConversationConfig{
						ReplyError: false,
					},
					Existing: ConversationConfig{
						ReplyError: false,
					},
				},
			},
			args: args{isExistingConversation: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Conversation: tt.fields.Conversation,
			}
			if got := c.ShouldReplyError(tt.args.isExistingConversation); got != tt.want {
				t.Errorf("ShouldReplyError() = %v, want %v", got, tt.want)
			}
		})
	}
}
