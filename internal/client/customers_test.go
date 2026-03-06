package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListCustomers(t *testing.T) {
	responseData := map[string]any{
		"customers": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "c1",
					"name":      "Acme Corp",
					"logoUrl":   "https://example.com/logo.png",
					"slugId":    "acme",
					"revenue":   100000,
					"size":      50.0,
					"createdAt": "2024-01-01T00:00:00Z",
					"owner":     map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
					"status":    map[string]any{"id": "s1", "name": "active", "displayName": "Active", "color": "#00ff00"},
					"tier":      map[string]any{"id": "t1", "name": "enterprise", "displayName": "Enterprise", "color": "#gold"},
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	customers, err := c.ListCustomers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(customers) != 1 {
		t.Fatalf("expected 1 customer, got %d", len(customers))
	}
	if customers[0].Name != "Acme Corp" {
		t.Errorf("expected Acme Corp, got %q", customers[0].Name)
	}
	if customers[0].Owner == nil || customers[0].Owner.Name != "Alice" {
		t.Errorf("expected owner Alice")
	}
	if customers[0].Status == nil || customers[0].Status.DisplayName != "Active" {
		t.Errorf("expected status Active")
	}
}

func TestListCustomersEmpty(t *testing.T) {
	responseData := map[string]any{
		"customers": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	customers, err := c.ListCustomers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(customers) != 0 {
		t.Fatalf("expected 0 customers, got %d", len(customers))
	}
}

func TestGetCustomer(t *testing.T) {
	responseData := map[string]any{
		"customer": map[string]any{
			"id":        "c1",
			"name":      "Acme Corp",
			"logoUrl":   "https://example.com/logo.png",
			"slugId":    "acme",
			"revenue":   100000,
			"size":      50.0,
			"createdAt": "2024-01-01T00:00:00Z",
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	customer, err := c.GetCustomer("c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.Name != "Acme Corp" {
		t.Errorf("expected Acme Corp, got %q", customer.Name)
	}
}

func TestGetCustomerNotFound(t *testing.T) {
	responseData := map[string]any{
		"customer": nil,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	_, err := c.GetCustomer("nonexistent")
	if err == nil {
		t.Fatal("expected error for not found customer")
	}
}

func TestListCustomerNeeds(t *testing.T) {
	responseData := map[string]any{
		"customerNeeds": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "cn1",
					"body":      "We need feature X",
					"priority":  1.0,
					"createdAt": "2024-01-01T00:00:00Z",
					"creator":   map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice", "email": "alice@test.com"},
					"customer":  map[string]any{"id": "c1", "name": "Acme Corp"},
					"issue":     map[string]any{"id": "i1", "identifier": "ENG-1", "title": "Feature X"},
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	needs, err := c.ListCustomerNeeds()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(needs) != 1 {
		t.Fatalf("expected 1 need, got %d", len(needs))
	}
	if needs[0].Body != "We need feature X" {
		t.Errorf("expected body 'We need feature X', got %q", needs[0].Body)
	}
	if needs[0].Customer == nil || needs[0].Customer.Name != "Acme Corp" {
		t.Errorf("expected customer Acme Corp")
	}
}

func TestCreateCustomer(t *testing.T) {
	responseData := map[string]any{
		"customerCreate": map[string]any{
			"success": true,
			"customer": map[string]any{
				"id":        "c1",
				"name":      "New Customer",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	customer, err := c.CreateCustomer("New Customer", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.Name != "New Customer" {
		t.Errorf("expected New Customer, got %q", customer.Name)
	}
}

func TestUpdateCustomer(t *testing.T) {
	responseData := map[string]any{
		"customerUpdate": map[string]any{
			"success": true,
			"customer": map[string]any{
				"id":        "c1",
				"name":      "Updated Customer",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	customer, err := c.UpdateCustomer("c1", map[string]any{"name": "Updated Customer"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.Name != "Updated Customer" {
		t.Errorf("expected Updated Customer, got %q", customer.Name)
	}
}

func TestDeleteCustomer(t *testing.T) {
	responseData := map[string]any{
		"customerDelete": map[string]any{
			"success": true,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	if err := c.DeleteCustomer("c1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpsertCustomer(t *testing.T) {
	responseData := map[string]any{
		"customerUpsert": map[string]any{
			"success": true,
			"customer": map[string]any{
				"id":        "c1",
				"name":      "Upserted Customer",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	customer, err := c.UpsertCustomer(map[string]any{"name": "Upserted Customer", "externalId": "ext-123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.Name != "Upserted Customer" {
		t.Errorf("expected Upserted Customer, got %q", customer.Name)
	}
}
