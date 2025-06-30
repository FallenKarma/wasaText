package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/fallenkarma/wasatext/internal/models"
	"github.com/fallenkarma/wasatext/internal/repository"
	"github.com/google/uuid"
)

// WASATextService implements Service
type WASATextService struct {
	repo repository.Repository
}

// NewWASATextService creates a new WASATextService
func NewWASATextService(repo repository.Repository) *WASATextService {
	return &WASATextService{
		repo: repo,
	}
}

// Login implements UserService.Login
func (s *WASATextService) Login(ctx context.Context, name string) (*models.User, error) {
	if len(name) < 3 || len(name) > 16 {
		return nil, errors.New("username must be between 3 and 16 characters")
	}

	return s.repo.CreateUser(ctx, name)
}

// GetUserByID implements UserService.GetUserByID
func (s *WASATextService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// GetAllUsers implements UserService.GetAllUsers
func (s *WASATextService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAllUsers(ctx)
}

// UpdateUsername implements UserService.UpdateUsername
func (s *WASATextService) UpdateUsername(ctx context.Context, userID string, newName string) error {
	if len(newName) < 3 || len(newName) > 16 {
		return errors.New("username must be between 3 and 16 characters")
	}
	
	return s.repo.UpdateUsername(ctx, userID, newName)
}

// SetUserPhoto implements UserService.SetUserPhoto
func (s *WASATextService) SetUserPhoto(ctx context.Context, userID string, photo multipart.File) (string, error) {
	return s.repo.SaveUserPhoto(ctx, userID, photo)
}

// CreateConversation creates a new conversation between users
func (s *Service) CreateConversation(ctx context.Context, creatorID string, participantIDs []string, Type models.ConversationType, Name string) (*models.Conversation, error) {
	// Validate participants exist
	for _, participantID := range participantIDs {
		if _, err := s.GetUser(ctx, participantID); err != nil {
			return nil, fmt.Errorf("invalid participant ID: %s", participantID)
		}
	}

	// Create list of all participants including the creator
	allParticipants := append([]string{creatorID}, participantIDs...)


	var conv *models.Conversation
	var err error
	if Type == models.DirectConversation {
		// Check if a DM conversation already exists between these two users
		conv, err = s.repo.CreateDirectConversation(ctx, creatorID, participantIDs[0])
		if err != nil {
			return nil, err
		}
	} else {
		conv, err = s.repo.CreateGroupConversation(ctx, Name, allParticipants)
		if err != nil {
			return nil, err
		}
	}


	return conv, nil
}



// GetConversation implements ConversationService.GetConversation
func (s *WASATextService) GetConversation(ctx context.Context, id string, userID string) (*models.Conversation, error) {
	// Get the conversation
	conv, err := s.repo.GetConversationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if conv == nil {
		return nil, errors.New("conversation not found")
	}

	// Check if user is a participant
	isParticipant := false
	for _, participant := range conv.Participants {
		if participant.ID == userID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant in this conversation")
	}

	// Mark all received messages as read
	for i, msg := range conv.Messages {
		if msg.Sender.ID != userID && msg.Status == models.Received {
			err := s.repo.UpdateMessageStatus(ctx, msg.ConversationID, models.Read)
			if err != nil {
				return nil, err
			}
			conv.Messages[i].Status = models.Read
		}
	}

	return conv, nil
}

// CreateDirectConversation implements ConversationService.CreateDirectConversation
func (s *WASATextService) CreateDirectConversation(ctx context.Context, userID1, userID2 string) (*models.Conversation, error) {
	// Validate users exist
	user1, err := s.repo.GetUserByID(ctx, userID1)
	if err != nil {
		return nil, err
	}
	if user1 == nil {
		return nil, errors.New("first user not found")
	}

	user2, err := s.repo.GetUserByID(ctx, userID2)
	if err != nil {
		return nil, err
	}
	if user2 == nil {
		return nil, errors.New("second user not found")
	}

	return s.repo.CreateDirectConversation(ctx, userID1, userID2)
}

// CreateGroupConversation implements ConversationService.CreateGroupConversation
func (s *WASATextService) CreateGroupConversation(ctx context.Context, name string, creatorID string, participants []string) (*models.Conversation, error) {
	// Validate group name
	if name == "" {
		return nil, errors.New("group name cannot be empty")
	}
	
	// Validate creator is included in participants
	creatorIncluded := false
	for _, participant := range participants {
		if participant == creatorID {
			creatorIncluded = true
			break
		}
	}
	if !creatorIncluded {
		participants = append(participants, creatorID)
	}
	
	// Validate all participants exist
	for _, participant := range participants {
		user, err := s.repo.GetUserByID(ctx, participant)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("one or more participants not found")
		}
	}
	
	return s.repo.CreateGroupConversation(ctx, name, participants)
}

// AddToGroup implements ConversationService.AddToGroup
func (s *WASATextService) AddToGroup(ctx context.Context, groupID, userID, currentUserID string) error {
	// Check if the current user is in the group
	conversation, err := s.repo.GetConversationByID(ctx, groupID)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.New("group not found")
	}
	
	// Verify it's a group conversation
	if conversation.Type != models.GroupConversation {
		return errors.New("cannot add users to a direct conversation")
	}
	
	// Check if the current user is a participant
	currentUserIsParticipant := false
	for _, participant := range conversation.Participants {
		if participant.ID == currentUserID {
			currentUserIsParticipant = true
			break
		}
	}
	if !currentUserIsParticipant {
		return errors.New("you are not a member of this group")
	}
	
	// Check if the user exists
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	return s.repo.AddUserToGroup(ctx, groupID, userID)
}

// LeaveGroup implements ConversationService.LeaveGroup
func (s *WASATextService) LeaveGroup(ctx context.Context, groupID, userID string) error {
	// Get the group
	conversation, err := s.repo.GetConversationByID(ctx, groupID)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.New("group not found")
	}
	
	// Verify it's a group conversation
	if conversation.Type != models.GroupConversation {
		return errors.New("cannot leave a direct conversation")
	}
	
	// Verify the user is in the group
	userIsParticipant := false
	for _, participant := range conversation.Participants {
		if participant.ID == userID {
			userIsParticipant = true
			break
		}
	}
	if !userIsParticipant {
		return errors.New("you are not a member of this group")
	}
	
	return s.repo.RemoveUserFromGroup(ctx, groupID, userID)
}

// SetGroupName implements ConversationService.SetGroupName
func (s *WASATextService) SetGroupName(ctx context.Context, groupID, name, userID string) error {
	// Validate group name
	if name == "" {
		return errors.New("group name cannot be empty")
	}
	
	// Check if the current user is in the group
	conversation, err := s.repo.GetConversationByID(ctx, groupID)
	if err != nil {
		return err
	}
	if conversation == nil {
		return errors.New("group not found")
	}
	
	// Verify it's a group conversation
	if conversation.Type != models.GroupConversation {
		return errors.New("cannot set name for a direct conversation")
	}
	
	// Check if the current user is a participant
	userIsParticipant := false
	for _, participant := range conversation.Participants {
		if participant.ID == userID {
			userIsParticipant = true
			break
		}
	}
	if !userIsParticipant {
		return errors.New("you are not a member of this group")
	}
	
	return s.repo.UpdateGroupName(ctx, groupID, name)
}

// SetGroupPhoto implements ConversationService.SetGroupPhoto
func (s *WASATextService) SetGroupPhoto(ctx context.Context, groupID string, photo multipart.File, userID string) (string, error) {
	// Check if the current user is in the group
	conversation, err := s.repo.GetConversationByID(ctx, groupID)
	if err != nil {
		return "", err
	}
	if conversation == nil {
		return "", errors.New("group not found")
	}
	
	// Verify it's a group conversation
	if conversation.Type != models.GroupConversation {
		return "", errors.New("cannot set photo for a direct conversation")
	}
	
	// Check if the current user is a participant
	userIsParticipant := false
	for _, participant := range conversation.Participants {
		if participant.ID == userID {
			userIsParticipant = true
			break
		}
	}
	if !userIsParticipant {
		return "", errors.New("you are not a member of this group")
	}
	
	return s.repo.SaveGroupPhoto(ctx, groupID, photo)
}

// SendMessage implements MessageService.SendMessage
func (s *WASATextService) SendMessage(ctx context.Context, conversationID, senderID, content string, messageType models.MessageType) (*models.Message, error) {
	// Validate the sender exists
	sender, err := s.repo.GetUserByID(ctx, senderID)
	if err != nil {
		return nil, err
	}
	if sender == nil {
		return nil, errors.New("sender not found")
	}
	
	// Validate the conversation exists
	conversation, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conversation == nil {
		return nil, errors.New("conversation not found")
	}
	
	// Validate the sender is a participant
	senderIsParticipant := false
	for _, participant := range conversation.Participants {
		if participant.ID == senderID {
			senderIsParticipant = true
			break
		}
	}
	if !senderIsParticipant {
		return nil, errors.New("sender is not a participant in this conversation")
	}
	
	// Create and send the message
	message := models.Message{
		ConversationID:        uuid.New().String(),
		Sender:    *sender,
		Timestamp: time.Now(),
		Content:   content,
		Type:      messageType,
		Status:    models.Sent,
	}
	
	createdMessage, err := s.repo.CreateMessage(ctx, message, conversationID)
	if err != nil {
		return nil, err
	}
	
	// Update message status to received for polling
	go func() {
		// Wait a bit to simulate network delay
		time.Sleep(1 * time.Second)
		
		// Update message status to received
		ctx := context.Background()
		err := s.repo.UpdateMessageStatus(ctx, createdMessage.ConversationID, models.Received)
		if err != nil {
			// Just log the error, don't propagate it
			// In a real app, we would use a logger
			// log.Printf("Error updating message status: %v", err)
		}
	}()
	
	return createdMessage, nil
}

// SendPhotoMessage implements MessageService.SendPhotoMessage
func (s *WASATextService) SendPhotoMessage(ctx context.Context, conversationID, senderID string, photo multipart.File) (*models.Message, error) {
	// Save the photo
	photoURL, err := s.repo.SaveMessagePhoto(ctx, senderID, photo)
	if err != nil {
		return nil, err
	}
	
	// Send a message with the photo URL
	return s.SendMessage(ctx, conversationID, senderID, photoURL, models.PhotoMessage)
}

// ForwardMessage implements MessageService.ForwardMessage
func (s *WASATextService) ForwardMessage(ctx context.Context, messageID, targetConversationID, userID string) error {
	// Get the message
	message, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.New("message not found")
	}
	sender, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if sender == nil {
		return errors.New("sender not found")
	}
	
	// Get the target conversation
	targetConversation, err := s.repo.GetConversationByID(ctx, targetConversationID)
	if err != nil {
		return err
	}
	if targetConversation == nil {
		return errors.New("target conversation not found")
	}
	
	// Check if the user is a participant in the target conversation
	userIsParticipant := false
	for _, participant := range targetConversation.Participants {
		if participant.ID == userID {
			userIsParticipant = true
			break
		}
	}
	if !userIsParticipant {
		return errors.New("user is not a participant in the target conversation")
	}
	
	// Create a new message in the target conversation
	newMessage := models.Message{
		ConversationID:        uuid.New().String(),
		Sender:    *sender,
		Timestamp: time.Now(),
		Content:   message.Content,
		Type:      message.Type,
		Status:    models.Sent,
	}
	
	_, err = s.repo.CreateMessage(ctx, newMessage, targetConversationID)
	return err
}

// DeleteMessage implements MessageService.DeleteMessage
func (s *WASATextService) DeleteMessage(ctx context.Context, messageID, userID string) error {
	// Get the message
	message, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.New("message not found")
	}
	
	// Check if the user is the sender of the message
	if message.Sender.ID != userID {
		return errors.New("only the sender can delete the message")
	}
	
	return s.repo.DeleteMessage(ctx, messageID)
}

// GetConversationMessages implements MessageService.GetConversationMessages
func (s *WASATextService) GetConversationMessages(ctx context.Context, conversationID string) ([]models.Message, error) {
	return s.repo.GetMessagesByConversationID(ctx, conversationID)
}

// MarkMessageAsReceived implements MessageService.MarkMessageAsReceived
func (s *WASATextService) MarkMessageAsReceived(ctx context.Context, messageID string) error {
	return s.repo.UpdateMessageStatus(ctx, messageID, models.Received)
}

// MarkMessageAsRead implements MessageService.MarkMessageAsRead
func (s *WASATextService) MarkMessageAsRead(ctx context.Context, messageID string) error {
	return s.repo.UpdateMessageStatus(ctx, messageID, models.Read)
}

// AddReaction implements MessageService.AddReaction
func (s *WASATextService) AddReaction(ctx context.Context, messageID, userID, emoji string) error {
	// Validate the message exists
	message, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return errors.New("message not found")
	}
	
	return s.repo.AddReaction(ctx, messageID, userID, emoji)
}

// RemoveReaction implements MessageService.RemoveReaction
func (s *WASATextService) RemoveReaction(ctx context.Context, messageID, userID string) error {
	return s.repo.RemoveReaction(ctx, messageID, userID)
}