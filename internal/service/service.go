package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/fallenkarma/wasatext/internal/models"
	"github.com/fallenkarma/wasatext/internal/repository"
)

// Service defines the business logic for the WASAText application
type Service struct {
	repo repository.Repository
}

// New creates a new service
func New(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Login authenticates a user or creates a new user if the username doesn't exist
func (s *Service) Login(ctx context.Context, username string) (*models.LoginResponse, error) {
	if len(username) < 3 || len(username) > 16 {
		return nil, errors.New("username must be between 3 and 16 characters")
	}

	user, err := s.repo.CreateUser(ctx, username)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Id: user.ID,
	}, nil
}

// UpdateUsername updates a user's username
func (s *Service) UpdateUsername(ctx context.Context, userID string, newUsername string) error {
	if len(newUsername) < 3 || len(newUsername) > 16 {
		return errors.New("username must be between 3 and 16 characters")
	}

	return s.repo.UpdateUsername(ctx, userID, newUsername)
}

// SetUserPhoto sets a user's profile photo
func (s *Service) SetUserPhoto(ctx context.Context, userID string, photo multipart.File) (string, error) {
	return s.repo.SaveUserPhoto(ctx, userID, photo)
}

// GetUser gets a user by ID
func (s *Service) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// GetUserByName gets a user by username
func (s *Service) GetUserByName(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetUserByName(ctx, username)
}

// GetAllUsers gets all users
func (s *Service) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAllUsers(ctx)
}

// GetConversations gets all conversations for a user
func (s *Service) GetConversations(ctx context.Context, userID string) ([]models.Conversation, error) {
	return s.repo.GetConversationsByUserID(ctx, userID)
}

// GetConversation gets a specific conversation
func (s *Service) GetConversation(ctx context.Context, conversationID string) (*models.Conversation, error) {
	return s.repo.GetConversationByID(ctx, conversationID)
}

// CreateDirectConversation creates a new direct conversation between two users
func (s *Service) CreateDirectConversation(ctx context.Context, userID1, userID2 string) (*models.Conversation, error) {
	// Validate users exist
	_, err := s.repo.GetUserByID(ctx, userID1)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.GetUserByID(ctx, userID2)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateDirectConversation(ctx, userID1, userID2)
}

// CreateGroupConversation creates a new group conversation
func (s *Service) CreateGroupConversation(ctx context.Context, name string, creatorID string, participants []string) (*models.Conversation, error) {
	// Make sure the creator is included in participants
	hasCreator := false
	for _, id := range participants {
		if id == creatorID {
			hasCreator = true
			break
		}
	}
	if !hasCreator {
		participants = append(participants, creatorID)
	}

	// Validate all participants exist
	for _, id := range participants {
		_, err := s.repo.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return s.repo.CreateGroupConversation(ctx, name, participants)
}

// AddToGroup adds a user to a group
func (s *Service) AddToGroup(ctx context.Context, groupID, userID string) error {
	// Check if the user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repo.AddUserToGroup(ctx, groupID, userID)
}

// LeaveGroup removes a user from a group
func (s *Service) LeaveGroup(ctx context.Context, groupID, userID string) error {
	return s.repo.RemoveUserFromGroup(ctx, groupID, userID)
}

// SetGroupName sets a group's name
func (s *Service) SetGroupName(ctx context.Context, groupID, name string) error {
	return s.repo.UpdateGroupName(ctx, groupID, name)
}

// SetGroupPhoto sets a group's photo
func (s *Service) SetGroupPhoto(ctx context.Context, groupID string, photo multipart.File) (string, error) {
	return s.repo.SaveGroupPhoto(ctx, groupID, photo)
}

// SendTextMessage sends a new text message
func (s *Service) SendTextMessage(ctx context.Context, senderID, conversationID, content string, replyToID *string) (*models.Message, error) {
	// Verify the conversation exists and the user is a participant
	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conv == nil {
		return nil, errors.New("conversation not found")
	}
	sender, err := s.repo.GetUserByID(ctx, senderID)
	if err != nil {
		return nil, err
	}
	if sender == nil {
		return nil, errors.New("sender not found")
	}

	// Check if the user is a participant in the conversation
	isParticipant := false
	for _, participant := range conv.Participants {
		if participant.ID == senderID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant in the conversation")
	}

	// Create the message
	msg := models.Message{
		Sender:    *sender,
		Content:   content,
		Type:      models.TextMessage,
		Status:    models.Sent,
	}

    if replyToID != nil && *replyToID != "" {
        msg.ReplyTo = replyToID
    }

	return s.repo.CreateMessage(ctx, msg, conversationID)
}

// SendPhotoMessage sends a new photo message
func (s *Service) SendPhotoMessage(ctx context.Context, senderID, conversationID string, photo multipart.File, replyToID string) (*models.Message, error) {
	// Verify the conversation exists and the user is a participant
	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conv == nil {
		return nil, errors.New("conversation not found")
	}
	sender, err := s.repo.GetUserByID(ctx, senderID)
	if err != nil {
		return nil, err
	}
	if sender == nil {
		return nil, errors.New("sender not found")
	}

	// Check if the user is a participant in the conversation
	isParticipant := false
	for _, participant := range conv.Participants {
		if participant.ID == senderID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant in the conversation")
	}

	// Save the photo and get the path
	photoPath, err := s.repo.SaveMessagePhoto(ctx, senderID, photo)
	if err != nil {
		return nil, err
	}

	// Create the message
	msg := models.Message{
		Sender:    *sender,
		Content:   photoPath,
		Type:      models.PhotoMessage,
		Status:    models.Sent,
	}


	if replyToID != "" {
		msg.ReplyTo = &replyToID
	}
	return s.repo.CreateMessage(ctx, msg, conversationID)
}

// ForwardMessage forwards a message to another conversation
func (s *Service) ForwardMessage(ctx context.Context, userID, messageID, targetConversationID string) error {
	// Get the original message
	msg, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("message not found")
	}

	// Verify the target conversation exists and the user is a participant
	targetConv, err := s.repo.GetConversationByID(ctx, targetConversationID)
	if err != nil {
		return err
	}
	if targetConv == nil {
		return errors.New("target conversation not found")
	}
	sender, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if sender == nil {
		return errors.New("sender not found")
	}

	// Check if the user is a participant in the target conversation
	isParticipant := false
	for _, participant := range targetConv.Participants {
		if participant.ID == userID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return errors.New("user is not a participant in the target conversation")
	}

	// Create a new message in the target conversation with the same content
	newMsg := models.Message{
		Sender:    *sender,
		Content:   msg.Content,
		Type:      msg.Type,
		Status:    models.Sent,
	}

	_, err = s.repo.CreateMessage(ctx, newMsg, targetConversationID)
	return err
}

// DeleteMessage deletes a message
func (s *Service) DeleteMessage(ctx context.Context, userID, messageID string) error {
	// Get the message
	msg, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("message not found")
	}
	

	// Check if the user is the sender of the message
	if msg.Sender.ID != userID {
		return errors.New("only the sender can delete a message")
	}

	return s.repo.DeleteMessage(ctx, messageID)
}

// UpdateMessage updates a message
func (s *Service) UpdateMessage(ctx context.Context, userID, messageID string, content string) error {
	// Get the message
	msg, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("message not found")
	}
	

	// Check if the user is the sender of the message
	if msg.Sender.ID != userID {
		return errors.New("only the sender can update a message")
	}

	return s.repo.UpdateMessageContent(ctx, messageID, content)
}

// AddReaction adds a reaction to a message
func (s *Service) AddReaction(ctx context.Context, userID, messageID, emoji string) error {
	// Get the message
	msg, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("message not found")
	}

	return s.repo.AddReaction(ctx, messageID, userID, emoji)
}

// RemoveReaction removes a reaction from a message
func (s *Service) RemoveReaction(ctx context.Context, userID, messageID string) error {
	return s.repo.RemoveReaction(ctx, messageID, userID)
}

// UpdateMessageStatus updates the status of a message
func (s *Service) UpdateMessageStatus(ctx context.Context, messageID string, status models.MessageStatus) error {
	return s.repo.UpdateMessageStatus(ctx, messageID, status)
}