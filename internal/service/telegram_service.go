package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func SendOTPTelegram(chatID, message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	botToken = strings.TrimSpace(botToken)
	botToken = strings.TrimPrefix(botToken, "bot")
	
	if botToken == "" {
		fmt.Printf("[MOCK TELEGRAM] MSG to chatID %s: %s\n", chatID, message)
		return nil
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML", // Mengubah parse mode menjadi HTML untuk memudahkan escape karakter
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\n[FALLBACK MSG] Failed to call Telegram API. MSG to chatID %s: %s\n\n", chatID, message)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("\n[FALLBACK MSG] Telegram API error (status %d). MSG to chatID %s: %s\n\n", resp.StatusCode, chatID, message)
		return fmt.Errorf("failed to send telegram msg: status %d", resp.StatusCode)
	}

	return nil
}
