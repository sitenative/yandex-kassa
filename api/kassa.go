package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dbzer0/yandex-kassa/api/client"
	"github.com/dbzer0/yandex-kassa/api/info"
	"github.com/dbzer0/yandex-kassa/api/payment"
)

const apiURL = "https://payment.yandex.net/api/v3"

type Kassa struct {
	MaxAttempts int
	client      *client.APIClient
}

// New создает объект для работы с API Яндекс Кассы.
func New(shopID, secretKey string) *Kassa {
	return &Kassa{
		client: &client.APIClient{
			HTTP:   http.DefaultClient,
			APIURL: apiURL,
			ShopID: shopID,
			Secret: secretKey,
		},
	}
}

// NewHTTPClient перезаписывает http клиент, определенный по-умолчанию.
func (k *Kassa) NewHTTPClient(client *http.Client) {
	k.client.HTTP = client
}

// NewPayment создает объект NewPayment. Используется для создания платежа.
func (k *Kassa) NewPayment(value, currency string) *payment.NewPayment {
	return &payment.NewPayment{
		APIClient: k.client,
		Amount: payment.Amount{
			Value:    value,
			Currency: currency,
		},
	}
}

// Payment создает создает объект Payment по которому доступны операции:
//   * получения информации о платеже;
//   * подтверждение платежа;
//   * отмена платежа;
func (k *Kassa) Payment(paymentID string) *info.Payment {
	return &info.Payment{
		APIClient: k.client,
		ID:        paymentID,
	}
}

// Find позволяет получить информацию о текущем состоянии платежа по
// его уникальному идентификатору.
func (k *Kassa) Find(ctx context.Context, paymentID string) (*info.Payment, error) {
	p := &info.Payment{
		ID:        paymentID,
		APIClient: k.client,
	}

	reply, err := p.APIClient.Find(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	defer reply.Close()

	if err := json.NewDecoder(reply).Decode(&p); err != nil {
		return nil, err
	}

	if p != nil && p.Type != nil && *p.Type == "error" && p.Description != nil {
		return p, errors.New(*p.Description)
	}

	return p, nil
}

// Capture подтверждает вашу готовность принять платеж.
func (k *Kassa) Capture(ctx context.Context, paymentID string) (*info.Payment, error) {
	return nil, nil
}

// Cancel отменяет платеж, находящийся в статусе waiting_for_capture.
func (k *Kassa) Cancel(ctx context.Context, idempKey, paymentID string) (*info.Payment, error) {
	p := &info.Payment{
		ID:        paymentID,
		APIClient: k.client,
	}

	reply, err := p.APIClient.Cancel(ctx, idempKey, paymentID)
	if err != nil {
		return nil, err
	}
	defer reply.Close()

	if err := json.NewDecoder(reply).Decode(&p); err != nil {
		return nil, err
	}

	if p != nil && p.Type != nil && *p.Type == "error" && p.Description != nil {
		return p, errors.New(*p.Description)
	}

	return p, nil
}
