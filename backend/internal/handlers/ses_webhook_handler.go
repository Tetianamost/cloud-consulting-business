package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// SESWebhookHandler handles SES webhook notifications for email delivery status
type SESWebhookHandler struct {
	sesService         interfaces.SESService
	emailEventRecorder interfaces.EmailEventRecorder
	logger             *logrus.Logger
}

// NewSESWebhookHandler creates a new SES webhook handler
func NewSESWebhookHandler(
	sesService interfaces.SESService,
	emailEventRecorder interfaces.EmailEventRecorder,
	logger *logrus.Logger,
) *SESWebhookHandler {
	return &SESWebhookHandler{
		sesService:         sesService,
		emailEventRecorder: emailEventRecorder,
		logger:             logger,
	}
}

// HandleSESNotification handles incoming SES notifications via webhook
func (h *SESWebhookHandler) HandleSESNotification(c *gin.Context) {
	h.logger.Info("Received SES webhook notification")

	// Parse the notification from the request body
	var notification interfaces.SESNotification
	if err := c.ShouldBindJSON(&notification); err != nil {
		h.logger.WithError(err).Error("Failed to parse SES notification")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid notification format",
		})
		return
	}

	// Log the notification details
	h.logger.WithFields(logrus.Fields{
		"message_id":        notification.MessageID,
		"notification_type": notification.NotificationType,
		"source":            notification.Source,
		"destinations":      notification.Destination,
	}).Info("Processing SES notification")

	// Process the notification using the SES service
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	result, err := h.sesService.ProcessSESNotification(ctx, &notification)
	if err != nil {
		h.logger.WithError(err).Error("Failed to process SES notification")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to process notification",
		})
		return
	}

	// Update email event status based on the notification
	err = h.updateEmailEventStatus(ctx, &notification, result)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update email event status")
		// Don't return error to SES - we processed the notification successfully
		// but failed to update our internal records
	}

	// Return success response to SES
	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"message_id":        result.MessageID,
		"notification_type": result.NotificationType,
		"status":            result.Status,
		"processed_at":      result.ProcessedAt,
	})

	h.logger.WithFields(logrus.Fields{
		"message_id":        result.MessageID,
		"notification_type": result.NotificationType,
		"status":            result.Status,
	}).Info("SES notification processed successfully")
}

// HandleSNSConfirmation handles SNS subscription confirmation for SES notifications
func (h *SESWebhookHandler) HandleSNSConfirmation(c *gin.Context) {
	h.logger.Info("Received SNS subscription confirmation")

	// Parse SNS message
	var snsMessage map[string]interface{}
	if err := c.ShouldBindJSON(&snsMessage); err != nil {
		h.logger.WithError(err).Error("Failed to parse SNS message")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid SNS message format",
		})
		return
	}

	// Check if this is a subscription confirmation
	messageType, exists := snsMessage["Type"]
	if !exists {
		h.logger.Error("SNS message missing Type field")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid SNS message - missing Type",
		})
		return
	}

	switch messageType {
	case "SubscriptionConfirmation":
		h.logger.Info("Processing SNS subscription confirmation")

		// Extract subscription URL
		subscribeURL, exists := snsMessage["SubscribeURL"]
		if !exists {
			h.logger.Error("SNS subscription confirmation missing SubscribeURL")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Missing SubscribeURL in confirmation",
			})
			return
		}

		h.logger.WithField("subscribe_url", subscribeURL).Info("SNS subscription URL received")

		// In a production environment, you would make an HTTP GET request to the SubscribeURL
		// to confirm the subscription. For now, we'll just log it.
		c.JSON(http.StatusOK, gin.H{
			"success":       true,
			"message":       "Subscription confirmation received",
			"subscribe_url": subscribeURL,
		})

	case "Notification":
		h.logger.Info("Processing SNS notification")

		// Extract the actual SES notification from the SNS message
		messageBody, exists := snsMessage["Message"]
		if !exists {
			h.logger.Error("SNS notification missing Message field")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Missing Message in SNS notification",
			})
			return
		}

		// Parse the SES notification from the SNS message body
		var sesNotification interfaces.SESNotification
		if err := json.Unmarshal([]byte(messageBody.(string)), &sesNotification); err != nil {
			h.logger.WithError(err).Error("Failed to parse SES notification from SNS message")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid SES notification in SNS message",
			})
			return
		}

		// Process the SES notification
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		result, err := h.sesService.ProcessSESNotification(ctx, &sesNotification)
		if err != nil {
			h.logger.WithError(err).Error("Failed to process SES notification from SNS")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to process SES notification",
			})
			return
		}

		// Update email event status
		err = h.updateEmailEventStatus(ctx, &sesNotification, result)
		if err != nil {
			h.logger.WithError(err).Error("Failed to update email event status from SNS")
		}

		c.JSON(http.StatusOK, gin.H{
			"success":           true,
			"message_id":        result.MessageID,
			"notification_type": result.NotificationType,
			"status":            result.Status,
		})

	default:
		h.logger.WithField("message_type", messageType).Warn("Unknown SNS message type")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Unknown SNS message type: %v", messageType),
		})
	}
}

// updateEmailEventStatus updates the email event status based on SES notification
func (h *SESWebhookHandler) updateEmailEventStatus(ctx context.Context, notification *interfaces.SESNotification, result *interfaces.SESNotificationResult) error {
	if notification.MessageID == "" {
		return fmt.Errorf("missing message ID in notification")
	}

	// Convert SES notification status to domain email event status
	var status domain.EmailEventStatus
	var deliveredAt *time.Time
	var errorMessage string

	switch notification.NotificationType {
	case "Delivery":
		status = domain.EmailStatusDelivered
		deliveredAt = &notification.Timestamp

	case "Bounce":
		status = domain.EmailStatusBounced
		if notification.Bounce != nil {
			errorMessage = fmt.Sprintf("Bounce: %s - %s",
				notification.Bounce.BounceType,
				notification.Bounce.BounceSubType)

			// Add details from bounced recipients
			if len(notification.Bounce.BouncedRecipients) > 0 {
				recipient := notification.Bounce.BouncedRecipients[0]
				if recipient.DiagnosticCode != "" {
					errorMessage += fmt.Sprintf(" (%s)", recipient.DiagnosticCode)
				}
			}
		}

	case "Complaint":
		status = domain.EmailStatusSpam
		if notification.Complaint != nil {
			errorMessage = fmt.Sprintf("Complaint: %s",
				notification.Complaint.ComplaintFeedbackType)
		}

	default:
		return fmt.Errorf("unknown notification type: %s", notification.NotificationType)
	}

	// Update the email event status
	err := h.emailEventRecorder.UpdateEmailStatus(ctx, notification.MessageID, status, deliveredAt, errorMessage)
	if err != nil {
		return fmt.Errorf("failed to update email event status: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"message_id":    notification.MessageID,
		"status":        status,
		"delivered_at":  deliveredAt,
		"error_message": errorMessage,
	}).Info("Email event status updated successfully")

	return nil
}

// GetWebhookStatus returns the status of the webhook handler
func (h *SESWebhookHandler) GetWebhookStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  "active",
		"message": "SES webhook handler is active and ready to receive notifications",
		"endpoints": map[string]string{
			"ses_notification": "/webhook/ses/notification",
			"sns_confirmation": "/webhook/ses/sns",
		},
	})
}
