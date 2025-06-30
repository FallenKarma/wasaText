package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/fallenkarma/wasatext/internal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// PostgresRepository implements the Repository interface
type PostgresRepository struct {
	db          *sql.DB
	uploadPath  string
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(connStr string, uploadPath string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	ex, err := os.Executable()
    if err != nil {
        return nil, fmt.Errorf("failed to get executable path: %w", err)
    }
    exPath := filepath.Dir(ex) // This gives you C:\progettiPersonali\wasaProject\bin (if you build to bin)

    // Define your desired uploads directory relative to the executable
    // Example: Create an "uploads" folder next to your executable
    uploadsDir := filepath.Join(exPath, "uploads")

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db:          db,
		uploadPath:  uploadsDir,
	}, nil
}

// Close closes the database connection
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

// CreateUser implements UserRepository.CreateUser
func (r *PostgresRepository) CreateUser(ctx context.Context, name string) (*models.User, error) {
	// Check if user with this name already exists
	existingUser, _ := r.GetUserByName(ctx, name)
	if existingUser != nil {
		return existingUser, nil
	}

	// Generate a unique ID
	id := uuid.New().String()

	// Insert the new user
	query := "INSERT INTO users (id, name) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, id, name)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:   id,
		Name: name,
	}, nil
}

// GetUserByID implements UserRepository.GetUserByID
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := "SELECT id, name, photo_url FROM users WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	var photoURL sql.NullString
	err := row.Scan(&user.ID, &user.Name, &photoURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}

	return &user, nil
}

// GetUserByName implements UserRepository.GetUserByName
func (r *PostgresRepository) GetUserByName(ctx context.Context, name string) (*models.User, error) {
	query := "SELECT id, name, photo_url FROM users WHERE name = $1"
	row := r.db.QueryRowContext(ctx, query, name)

	var user models.User
	var photoURL sql.NullString
	err := row.Scan(&user.ID, &user.Name, &photoURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}

	return &user, nil
}

// UpdateUsername implements UserRepository.UpdateUsername
func (r *PostgresRepository) UpdateUsername(ctx context.Context, userID string, newName string) error {
	// Check if name is already in use
	existingUser, _ := r.GetUserByName(ctx, newName)
	if existingUser != nil && existingUser.ID != userID {
		return errors.New("username already in use")
	}

	query := "UPDATE users SET name = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, newName, userID)
	return err
}

// SaveUserPhoto implements UserRepository.SaveUserPhoto
func (r *PostgresRepository) SaveUserPhoto(ctx context.Context, userID string, photo multipart.File) (string, error) {
	// Create user photos directory if it doesn't exist
	userPhotosDir := filepath.Join(r.uploadPath, "user_photos")
	log.Println("User photos directory:", userPhotosDir)
	if err := os.MkdirAll(userPhotosDir, 0755); err != nil {
		return "", err
	}

	// Generate a unique filename
	filename := fmt.Sprintf("%s_%d.jpg", userID, time.Now().Unix())
	filepath := filepath.Join(userPhotosDir, filename)

	// Save the file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, photo); err != nil {
		return "", err
	}

	// Update the user's photo URL in the database
	relativePath := fmt.Sprintf("/uploads/user_photos/%s", filename)
	query := "UPDATE users SET photo_url = $1 WHERE id = $2"
	_, err = r.db.ExecContext(ctx, query, relativePath, userID)
	if err != nil {
		return "", err
	}

	return relativePath, nil
}

// GetAllUsers implements UserRepository.GetAllUsers
func (r *PostgresRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := "SELECT id, name, photo_url FROM users"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var photoURL sql.NullString
		if err := rows.Scan(&user.ID, &user.Name, &photoURL); err != nil {
			return nil, err
		}
		if photoURL.Valid {
			user.PhotoURL = photoURL.String
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// CreateDirectConversation implements ConversationRepository.CreateDirectConversation
func (r *PostgresRepository) CreateDirectConversation(ctx context.Context, userID1, userID2 string) (*models.Conversation, error) {
	// Check if a direct conversation already exists between these users
	log.Print("Users", userID1, userID2)
	query := `
		    SELECT 
        		conversation_id
			FROM 
				conversation_participants
			JOIN conversations ON conversation_participants.conversation_id = conversations.id
			WHERE 
				user_id IN ($1,$2) AND conversations.type = 'direct'
			GROUP BY 
				conversation_id
			HAVING 
				COUNT(DISTINCT user_id) = 2
	`
	row := r.db.QueryRowContext(ctx, query, userID1, userID2)

	var conversationID string
	err := row.Scan(&conversationID)
	if err == nil {
		// Conversation exists, return it
		return r.GetConversationByID(ctx, conversationID)
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Unexpected error
		return nil, err
	}

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create a new conversation
	id := uuid.New().String()
	insertConvQuery := "INSERT INTO conversations (id, type) VALUES ($1, $2)"
	_, err = tx.ExecContext(ctx, insertConvQuery, id, models.DirectConversation)
	if err != nil {
		return nil, err
	}

	// Add participants
	insertPartQuery := "INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1, $2)"
	_, err = tx.ExecContext(ctx, insertPartQuery, id, userID1)
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, insertPartQuery, id, userID2)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetConversationByID(ctx, id)
}

// CreateGroupConversation implements ConversationRepository.CreateGroupConversation
func (r *PostgresRepository) CreateGroupConversation(ctx context.Context, name string, participants []string) (*models.Conversation, error) {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create a new conversation
	id := uuid.New().String()
	insertConvQuery := "INSERT INTO conversations (id, name, type) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, insertConvQuery, id, name, models.GroupConversation)
	if err != nil {
		return nil, err
	}

	// Add participants
	insertPartQuery := "INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1, $2)"
	for _, userID := range participants {
		_, err = tx.ExecContext(ctx, insertPartQuery, id, userID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetConversationByID(ctx, id)
}

// GetConversationByID implements ConversationRepository.GetConversationByID
func (r *PostgresRepository) GetConversationByID(ctx context.Context, id string) (*models.Conversation, error) {
	// Get conversation details
	convQuery := "SELECT id, name, type, photo_url FROM conversations WHERE id = $1"
	convRow := r.db.QueryRowContext(ctx, convQuery, id)

	var conv models.Conversation
	var name, photoURL sql.NullString
	var convType string
	err := convRow.Scan(&conv.ID, &name, &convType, &photoURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if name.Valid {
		conv.Name = name.String
	}
	if photoURL.Valid {
		conv.PhotoURL = photoURL.String
	}
	conv.Type = models.ConversationType(convType)

	// Get participants
	partQuery := "SELECT cp.user_id, u.name, u.photo_url FROM conversation_participants cp JOIN users u ON u.id = cp.user_id  WHERE cp.conversation_id = $1"
	partRows, err := r.db.QueryContext(ctx, partQuery, id)
	if err != nil {
		return nil, err
	}
	defer partRows.Close()

	for partRows.Next() {
		var userID string
		var userName string
		var photo_url sql.NullString
		if err := partRows.Scan(&userID,&userName, &photo_url); err != nil {
			return nil, err
		}
		
		userPhotoUrl := ""
		if photo_url.Valid {
			userPhotoUrl = photo_url.String
		}
		
		conv.Participants = append(conv.Participants, models.Participant{
            ID:   userID,
            Name: userName,
			PhotoURL: userPhotoUrl,
        })

	}
	if err := partRows.Err(); err != nil {
		return nil, err
	}

	// Get messages
	messages, err := r.GetMessagesByConversationID(ctx, id)
	if err != nil {
		return nil, err
	}
	conv.Messages = messages

	// Set the last message if there are any messages
	if len(messages) > 0 {
		lastMsg := messages[len(messages)-1]  // Assuming messages are ordered by timestamp desc
		conv.LastMessage = &lastMsg
	}

	// If this is a direct conversation and has no name, set the name to the other user's name
	if conv.Type == models.DirectConversation && !name.Valid && len(conv.Participants) == 2 {
		// Find the other user in the conversation
		var otherUserID string
		for _, participant := range conv.Participants {
			if participant.ID != ctx.Value("userID").(string) {
				otherUserID = participant.ID
				break
			}
		}

		// Get the other user's name
		if otherUserID != "" {
			otherUser, err := r.GetUserByID(ctx, otherUserID)
			if err == nil && otherUser != nil {
				conv.Name = otherUser.Name
			}
		}
	}

	return &conv, nil
}

// GetConversationsByUserID implements ConversationRepository.GetConversationsByUserID
func (r *PostgresRepository) GetConversationsByUserID(ctx context.Context, userID string) ([]models.Conversation, error) {
	// Find all conversations where the user is a participant
	query := `
		SELECT c.id 
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1
		ORDER BY c.last_activity DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var convID string
		if err := rows.Scan(&convID); err != nil {
			return nil, err
		}

		conv, err := r.GetConversationByID(ctx, convID)
		if err != nil {
			return nil, err
		}
		if conv != nil {
			conversations = append(conversations, *conv)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// AddUserToGroup implements ConversationRepository.AddUserToGroup
func (r *PostgresRepository) AddUserToGroup(ctx context.Context, groupID, userID string) error {
	// Check if the conversation is a group
	convQuery := "SELECT type FROM conversations WHERE id = $1"
	var convType string
	err := r.db.QueryRowContext(ctx, convQuery, groupID).Scan(&convType)
	if err != nil {
		return err
	}
	if convType != string(models.GroupConversation) {
		return errors.New("conversation is not a group")
	}

	// Check if user is already in the group
	checkQuery := "SELECT COUNT(*) FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2"
	var count int
	err = r.db.QueryRowContext(ctx, checkQuery, groupID, userID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user is already in the group")
	}

	// Add user to the group
	insertQuery := "INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1, $2)"
	_, err = r.db.ExecContext(ctx, insertQuery, groupID, userID)
	return err
}

// RemoveUserFromGroup implements ConversationRepository.RemoveUserFromGroup
func (r *PostgresRepository) RemoveUserFromGroup(ctx context.Context, groupID, userID string) error {
	// Check if the conversation is a group
	convQuery := "SELECT type FROM conversations WHERE id = $1"
	var convType string
	err := r.db.QueryRowContext(ctx, convQuery, groupID).Scan(&convType)
	if err != nil {
		return err
	}
	if convType != string(models.GroupConversation) {
		return errors.New("conversation is not a group")
	}

	// Remove user from the group
	deleteQuery := "DELETE FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2"
	result, err := r.db.ExecContext(ctx, deleteQuery, groupID, userID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user is not in the group")
	}

	return nil
}

// UpdateGroupName implements ConversationRepository.UpdateGroupName
func (r *PostgresRepository) UpdateGroupName(ctx context.Context, groupID, name string) error {
	// Check if the conversation is a group
	convQuery := "SELECT type FROM conversations WHERE id = $1"
	var convType string
	err := r.db.QueryRowContext(ctx, convQuery, groupID).Scan(&convType)
	if err != nil {
		return err
	}
	if convType != string(models.GroupConversation) {
		return errors.New("conversation is not a group")
	}

	// Update the group name
	updateQuery := "UPDATE conversations SET name = $1 WHERE id = $2"
	_, err = r.db.ExecContext(ctx, updateQuery, name, groupID)
	return err
}

// SaveGroupPhoto implements ConversationRepository.SaveGroupPhoto
func (r *PostgresRepository) SaveGroupPhoto(ctx context.Context, groupID string, photo multipart.File) (string, error) {
	// Check if the conversation is a group
	convQuery := "SELECT type FROM conversations WHERE id = $1"
	var convType string
	err := r.db.QueryRowContext(ctx, convQuery, groupID).Scan(&convType)
	if err != nil {
		return "", err
	}
	if convType != string(models.GroupConversation) {
		return "", errors.New("conversation is not a group")
	}

	// Create group photos directory if it doesn't exist
	groupPhotosDir := filepath.Join(r.uploadPath, "group_photos")
	if err := os.MkdirAll(groupPhotosDir, 0755); err != nil {
		return "", err
	}

	// Generate a unique filename
	filename := fmt.Sprintf("%s_%d.jpg", groupID, time.Now().Unix())
	filepath := filepath.Join(groupPhotosDir, filename)

	// Save the file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, photo); err != nil {
		return "", err
	}

	// Update the group's photo URL in the database
	relativePath := fmt.Sprintf("/uploads/group_photos/%s", filename)
	query := "UPDATE conversations SET photo_url = $1 WHERE id = $2"
	_, err = r.db.ExecContext(ctx, query, relativePath, groupID)
	if err != nil {
		return "", err
	}

	return relativePath, nil
}

// CreateMessage implements MessageRepository.CreateMessage
func (r *PostgresRepository) CreateMessage(ctx context.Context, msg models.Message, conversationID string) (*models.Message, error) {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// If no ID provided, generate one
	if msg.ConversationID == "" {
		msg.ConversationID = uuid.New().String()
	}

	// If no timestamp provided, use current time
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	// Insert the message
	msgQuery := `
		INSERT INTO messages (id, sender_id, conversation_id, content, type, status, reply_to, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, msgQuery, msg.ConversationID, msg.Sender.ID, conversationID, msg.Content, msg.Type, msg.Status, msg.ReplyTo, msg.Timestamp)
	if err != nil {
		return nil, err
	}

	// Update the last activity timestamp of the conversation
	updateConvQuery := "UPDATE conversations SET last_activity = $1 WHERE id = $2"
	_, err = tx.ExecContext(ctx, updateConvQuery, msg.Timestamp, conversationID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &msg, nil
}

// GetMessagesByConversationID implements MessageRepository.GetMessagesByConversationID
func (r *PostgresRepository) GetMessagesByConversationID(ctx context.Context, conversationID string) ([]models.Message, error) {
	// Get messages with user information
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, u.name, u.photo_url, m.content, m.type, m.status, m.reply_to, m.timestamp, m.deleted_at
		FROM messages m
		INNER JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = $1
		ORDER BY m.timestamp ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var photoURL sql.NullString // Handle potential NULL photo_url

		
		if err := rows.Scan(   
			&msg.ID,                    // m.id              
			&msg.ConversationID,        // m.conversation_id
			&msg.Sender.ID,             // m.sender_id (User.ID)
			&msg.Sender.Name,           // u.name (User.Name)
			&photoURL,                  // u.photo_url (User.PhotoURL)
			&msg.Content,               // m.content
			&msg.Type,                  // m.type
			&msg.Status,                // m.status
			&msg.ReplyTo,                   // m.reply_to
			&msg.Timestamp,             // m.timestamp
			&msg.DeletedAt,             // m.deleted_at
		); err != nil {
			return nil, err
		}

		// Handle nullable photo URL
		if photoURL.Valid {
			msg.Sender.PhotoURL = photoURL.String
		}

		reactions, err := r.GetReactionsByMessageID(ctx, msg.ID)
		if err != nil {
			return nil, err
		}
		msg.Reactions = reactions


		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
// GetMessageByID implements MessageRepository.GetMessageByID
func (r *PostgresRepository) GetMessageByID(ctx context.Context, id string) (*models.Message, error) {
	query := `
		SELECT m.id, m.sender_id, u.name, m.content, m.type, m.status, m.reply_to, m.timestamp, m.conversation_id
		FROM messages m
		INNER JOIN users u ON m.sender_id = u.id
		WHERE m.id = $1 
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var msg models.Message
	var conversationID sql.NullString
	err := row.Scan(&msg.ConversationID, &msg.Sender.ID, &msg.Sender.Name, &msg.Content, &msg.Type, &msg.Status, &msg.ReplyTo, &msg.Timestamp, &conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &msg, nil
}

// DeleteMessage implements MessageRepository.DeleteMessage
func (r *PostgresRepository) DeleteMessage(ctx context.Context, id string) error {
	// Soft delete by setting the deleted_at timestamp
	query := "UPDATE messages SET deleted_at = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// UpdateMessageStatus implements MessageRepository.UpdateMessageStatus
func (r *PostgresRepository) UpdateMessageStatus(ctx context.Context, id string, status models.MessageStatus) error {
	query := "UPDATE messages SET status = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// UpdateMessageStatus implements MessageRepository.UpdateMessageStatus
func (r *PostgresRepository) UpdateMessageContent(ctx context.Context, id string, content string) error {
	query := "UPDATE messages SET content = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, content, id)
	return err
}


// SaveMessagePhoto implements MessageRepository.SaveMessagePhoto
func (r *PostgresRepository) SaveMessagePhoto(ctx context.Context, senderID string, photo multipart.File) (string, error) {
	// Create message photos directory if it doesn't exist
	msgPhotosDir := filepath.Join(r.uploadPath, "message_photos")
	if err := os.MkdirAll(msgPhotosDir, 0755); err != nil {
		return "", err
	}

	// Generate a unique filename
	filename := fmt.Sprintf("%s_%d.jpg", senderID, time.Now().Unix())
	filepath := filepath.Join(msgPhotosDir, filename)

	// Save the file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, photo); err != nil {
		return "", err
	}

	// Return the relative path to the file
	relativePath := fmt.Sprintf("/uploads/message_photos/%s", filename)
	return relativePath, nil
}

// AddReaction implements ReactionRepository.AddReaction
func (r *PostgresRepository) AddReaction(ctx context.Context, messageID, userID, emoji string) error {
	// Check if reaction already exists
	checkQuery := "SELECT COUNT(*) FROM reactions WHERE message_id = $1 AND user_id = $2"
	var count int
	err := r.db.QueryRowContext(ctx, checkQuery, messageID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Update existing reaction
		updateQuery := "UPDATE reactions SET emoji = $1 WHERE message_id = $2 AND user_id = $3"
		_, err = r.db.ExecContext(ctx, updateQuery, emoji, messageID, userID)
		return err
	}

	// Insert new reaction
	insertQuery := "INSERT INTO reactions (message_id, user_id, emoji) VALUES ($1, $2, $3)"
	_, err = r.db.ExecContext(ctx, insertQuery, messageID, userID, emoji)
	return err
}

// RemoveReaction implements ReactionRepository.RemoveReaction
func (r *PostgresRepository) RemoveReaction(ctx context.Context, messageID, userID string) error {
	deleteQuery := "DELETE FROM reactions WHERE message_id = $1 AND user_id = $2"
	result, err := r.db.ExecContext(ctx, deleteQuery, messageID, userID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("reaction not found")
	}

	return nil
}

// GetReactionsByMessageID implements ReactionRepository.GetReactionsByMessageID
func (r *PostgresRepository) GetReactionsByMessageID(ctx context.Context, messageID string) ([]models.Reaction, error) {
	query := "SELECT message_id, user_id, emoji FROM reactions WHERE message_id = $1"
	rows, err := r.db.QueryContext(ctx, query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []models.Reaction
	for rows.Next() {
		var reaction models.Reaction
		if err := rows.Scan(&reaction.MessageID, &reaction.UserID, &reaction.Emoji); err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}