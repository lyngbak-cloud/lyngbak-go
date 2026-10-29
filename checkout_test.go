package lyngbak

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestCreateResourceCheckout(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"message":"Please complete the checkout to create your resource.","checkout_url":"https://example.test/checkout/abc"}`)
	})
	res, err := client.Resources.Create(context.Background(), &CreateResourceRequest{PlanType: "pro"})
	if res != nil {
		t.Errorf("res = %+v, want nil", res)
	}
	var ce *CheckoutRequiredError
	if !errors.As(err, &ce) {
		t.Fatalf("err = %v, want *CheckoutRequiredError", err)
	}
	if ce.CheckoutURL != "https://example.test/checkout/abc" {
		t.Errorf("CheckoutURL = %q", ce.CheckoutURL)
	}
	if !strings.Contains(err.Error(), "https://example.test/checkout/abc") {
		t.Errorf("error %q does not contain the URL", err.Error())
	}
}

func TestUpdateResourcePaymentPending(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		fmt.Fprint(w, `{"message":"Resource update is awaiting payment confirmation.","order_id":"ord-1"}`)
	})
	res, err := client.Resources.Update(context.Background(), "r-1", &UpdateResourceRequest{})
	if res != nil {
		t.Errorf("res = %+v, want nil", res)
	}
	var pe *PaymentPendingError
	if !errors.As(err, &pe) {
		t.Fatalf("err = %v, want *PaymentPendingError", err)
	}
	if pe.OrderID != "ord-1" {
		t.Errorf("OrderID = %q", pe.OrderID)
	}
	if !strings.Contains(err.Error(), "ord-1") {
		t.Errorf("error %q does not contain the order ID", err.Error())
	}
}

func TestCreateResourceReturnsResource(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"uuid":"r-1","state":"pending"}`)
	})
	res, err := client.Resources.Create(context.Background(), &CreateResourceRequest{PlanType: "dev"})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil || res.UUID != "r-1" {
		t.Errorf("res = %+v, want UUID r-1", res)
	}
}

func TestResourceEmptyUUIDIsError(t *testing.T) {
	for name, body := range map[string]string{"empty object": `{}`, "message only": `{"message":"ok"}`} {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, body)
		})
		res, err := client.Resources.Create(context.Background(), &CreateResourceRequest{PlanType: "dev"})
		if err == nil || res != nil {
			t.Errorf("%s: Create res=%+v err=%v, want error", name, res, err)
		}
		res, err = client.Resources.Update(context.Background(), "r-1", &UpdateResourceRequest{})
		if err == nil || res != nil {
			t.Errorf("%s: Update res=%+v err=%v, want error", name, res, err)
		}
	}
}
