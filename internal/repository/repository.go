package repository

import (
	"context"
	"mime/multipart"

	"github.com/fallenkarma/wasatext/internal/models"
)

// UserRepository defines operations for user management
type UserRepository interface {
	// CreateUser creates a new user with the given name
	CreateUser(ctx context.Context, name string) (*models.User, error)
	
	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	
	// GetUserByName retrieves a user by their name
	GetUserByName(ctx context.Context, name string) (*models.User, error)
	
	// UpdateUsername updates a user's name
	UpdateUsername(ctx context.Context, userID string, newName string) error
	
	// SaveUserPhoto saves a user's profile photo
	SaveUserPhoto(ctx context.Context, userID string, photo multipart.File) (string, error)
	
	// GetAllUsers retrieves all users
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

// ConversationRepository defines operations for conversation management
type ConversationRepository interface {
	// CreateDirectConversation creates a new direct conversation between two users
	CreateDirectConversation(ctx context.Context, userID1, userID2 string) (*models.Conversation, error)
	
	// CreateGroupConversation creates a new group conversation
	CreateGroupConversation(ctx context.Context, name string, participants []string) (*models.Conversation, error)
	
	// GetConversationByID retrieves a conversation by its ID
	GetConversationByID(ctx context.Context, id string) (*models.Conversation, error)
	
	// GetConversationsByUserID retrieves all conversations for a user
	GetConversationsByUserID(ctx context.Context, userID string) ([]models.Conversation, error)
	
	// AddUserToGroup adds a user to a group conversation
	AddUserToGroup(ctx context.Context, groupID, userID string) error
	
	// RemoveUserFromGroup removes a user from a group conversation
	RemoveUserFromGroup(ctx context.Context, groupID, userID string) error
	
	// UpdateGroupName updates a group's name
	UpdateGroupName(ctx context.Context, groupID, name string) error
	
	// SaveGroupPhoto saves a group's photo
	SaveGroupPhoto(ctx context.Context, groupID string, photo multipart.File) (string, error)
}

// MessageRepository defines operations for message management
type MessageRepository interface {
	// CreateMessage creates a new message
	CreateMessage(ctx context.Context, msg models.Message, conversationID string) (*models.Message, error)
	
	// GetMessagesByConversationID retrieves all messages for a conversation
	GetMessagesByConversationID(ctx context.Context, conversationID string) ([]models.Message, error)
	
	// GetMessageByID retrieves a message by its ID
	GetMessageByID(ctx context.Context, id string) (*models.Message, error)
	
	// DeleteMessage marks a message as deleted
	DeleteMessage(ctx context.Context, id string) error
	
	// UpdateMessageStatus updates the status of a message
	UpdateMessageStatus(ctx context.Context, id string, status models.MessageStatus) error

	UpdateMessageContent(ctx context.Context, id string, content string) error
	
	// SaveMessagePhoto saves a photo message
	SaveMessagePhoto(ctx context.Context, senderID string, photo multipart.File) (string, error)
}

// ReactionRepository defines operations for reaction management
type ReactionRepository interface {
	// AddReaction adds a reaction to a message
	AddReaction(ctx context.Context, messageID, userID, emoji string) error
	
	// RemoveReaction removes a reaction from a message
	RemoveReaction(ctx context.Context, messageID, userID string) error
	
	// GetReactionsByMessageID retrieves all reactions for a message
	GetReactionsByMessageID(ctx context.Context, messageID string) ([]models.Reaction, error)
}

// Repository combines all repository interfaces
type Repository interface {
	UserRepository
	ConversationRepository
	MessageRepository
	ReactionRepository
}