package inter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Pix struct {
	client *Client
	token  Token
}

func NewPix(client *Client, token Token) *Pix {
	return &Pix{
		client: client,
		token:  token,
	}
}

type Calendar struct {
	Expirate int `json:"expiracao"`
}

type CalendarWishCreated struct {
	Expirate int       `json:"expiracao"`
	Created  time.Time `json:"criacao"`
}

type Debtor struct {
	CPF  string `json:"cpf,omitempty"`
	CNPJ string `json:"cnpj,omitempty"`
	Name string `json:"nome"`
}

type Value struct {
	Original  string `json:"original"`
	MobChange int    `json:"modalidadeAlteracao"`
}

type AdditionalInformation struct {
	Name  string `json:"nome"`
	Value string `json:"valor"`
}

type Loc struct {
	ID       int    `json:"id"`
	Location string `json:"location"`
	TypeCob  string `json:"tipoCob"`
}

type Charge struct {
	Calendar       Calendar                `json:"calendario"`
	Debtor         Debtor                  `json:"devedor"`
	Value          Value                   `json:"valor"`
	Key            string                  `json:"chave"`
	RequestText    string                  `json:"solicitacaoPagador"`
	AdditionalInfo []AdditionalInformation `json:"infoAdicionais"`
}

type ResponseCarge struct {
	CalendarWishCreated Calendar                `json:"calendario"`
	Txid                string                  `json:"txid"`
	Review              int                  `json:"revisao"`
	Loc                 Loc                     `json:"loc"`
	Status              string                  `json:"status"`
	Debtor              Debtor                  `json:"devedor"`
	Value               Value                   `json:"valor"`
	Key                 string                  `json:"chave"`
	Location            string                  `json:"location"`
	PixCopyAndPaste     string                  `json:"pixCopiaECola"`
	RequestText         string                  `json:"solicitacaoPagador"`
	AdditionalInfo      []AdditionalInformation `json:"infoAdicionais"`
}

func (p *Pix) NewChargeWithTxid(ctx context.Context, txid string, charge Charge) (ResponseCarge, error) {
	endpoint := fmt.Sprintf("%s/pix/v2/cob/%s", p.client.apiBaseUrl, txid)

	jsonData, _ := json.Marshal(charge)

	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return ResponseCarge{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token.Data))

	resp, err := p.client.Do(req)
	if err != nil {
		return ResponseCarge{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseCarge{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ResponseCarge{}, errors.New(string(data))
	}

	var tmp ResponseCarge

	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		return ResponseCarge{}, err
	}

	return tmp, nil
}

type ReviewCharge struct {
	Loc         Loc    `json:"loc"`
	Debtor      Debtor `json:"devedor"`
	Value       Value  `json:"valor"`
	RequestText string `json:"solicitacaoPagador"`
}

func (p *Pix) ReviewCharge(ctx context.Context, txid string, charge ReviewCharge) (ResponseCarge, error) {

	endpoint := fmt.Sprintf("%s/pix/v2/cob/%s", p.client.apiBaseUrl, txid)

	jsonData, _ := json.Marshal(charge)

	req, err := http.NewRequestWithContext(ctx, "PATCH", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return ResponseCarge{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token.Data))

	resp, err := p.client.Do(req)
	if err != nil {
		return ResponseCarge{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseCarge{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ResponseCarge{}, errors.New(string(data))
	}

	var tmp ResponseCarge

	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		return ResponseCarge{}, err
	}

	return tmp, nil
}

func (p *Pix) GetCharge(ctx context.Context, txid string) (ResponseCarge, error) {
	endpoint := fmt.Sprintf("%s/pix/v2/cob/%s", p.client.apiBaseUrl, txid)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return ResponseCarge{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token.Data))

	resp, err := p.client.Do(req)
	if err != nil {
		return ResponseCarge{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseCarge{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ResponseCarge{}, errors.New(string(data))
	}

	var tmp ResponseCarge

	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		return ResponseCarge{}, err
	}

	return tmp, nil
}

func (p *Pix) NewCharge(ctx context.Context, charge Charge) (ResponseCarge, error) {
	endpoint := fmt.Sprintf("%s/pix/v2/cob", p.client.apiBaseUrl)

	fmt.Printf("%s \n\n", endpoint)

	jsonData, _ := json.Marshal(charge)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return ResponseCarge{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.token.Data))

	resp, err := p.client.Do(req)
	if err != nil {
		return ResponseCarge{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseCarge{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return ResponseCarge{}, errors.New(string(data))
	}

	var tmp ResponseCarge


	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		return ResponseCarge{}, err
	}

	return tmp, nil

}
