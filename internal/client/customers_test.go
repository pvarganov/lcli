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

func TestListCustomerStatuses(t *testing.T) {
	responseData := map[string]any{
		"customerStatuses": map[string]any{
			"nodes": []map[string]any{
				{"id": "s1", "name": "active", "displayName": "Active", "color": "#00ff00", "description": "Active customer"},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	statuses, err := c.ListCustomerStatuses()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if statuses[0].DisplayName != "Active" {
		t.Errorf("expected Active, got %q", statuses[0].DisplayName)
	}
}

func TestListCustomerTiers(t *testing.T) {
	responseData := map[string]any{
		"customerTiers": map[string]any{
			"nodes": []map[string]any{
				{"id": "t1", "name": "enterprise", "displayName": "Enterprise", "color": "#gold", "description": "Enterprise tier"},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	tiers, err := c.ListCustomerTiers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tiers) != 1 {
		t.Fatalf("expected 1 tier, got %d", len(tiers))
	}
	if tiers[0].DisplayName != "Enterprise" {
		t.Errorf("expected Enterprise, got %q", tiers[0].DisplayName)
	}
}

func TestCreateCustomerNeed(t *testing.T) {
	responseData := map[string]any{
		"customerNeedCreate": map[string]any{
			"success": true,
			"customerNeed": map[string]any{
				"id":        "cn1",
				"body":      "We need feature X",
				"priority":  1.0,
				"createdAt": "2024-01-01T00:00:00Z",
				"customer":  map[string]any{"id": "c1", "name": "Acme Corp"},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	need, err := c.CreateCustomerNeed(map[string]any{"body": "We need feature X", "customerId": "c1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if need.Body != "We need feature X" {
		t.Errorf("expected body 'We need feature X', got %q", need.Body)
	}
}

func TestCreateCustomerNeedFailure(t *testing.T) {
	responseData := map[string]any{
		"customerNeedCreate": map[string]any{
			"success": false,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	_, err := c.CreateCustomerNeed(map[string]any{"body": "test"})
	if err == nil {
		t.Fatal("expected error on failure")
	}
}

func TestUpdateCustomerNeed(t *testing.T) {
	responseData := map[string]any{
		"customerNeedUpdate": map[string]any{
			"success": true,
			"customerNeed": map[string]any{
				"id":        "cn1",
				"body":      "Updated need",
				"priority":  2.0,
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
	need, err := c.UpdateCustomerNeed("cn1", map[string]any{"body": "Updated need"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if need.Body != "Updated need" {
		t.Errorf("expected 'Updated need', got %q", need.Body)
	}
}

func TestDeleteCustomerNeed(t *testing.T) {
	responseData := map[string]any{
		"customerNeedDelete": map[string]any{
			"success": true,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	if err := c.DeleteCustomerNeed("cn1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateCustomerStatus(t *testing.T) {
	responseData := map[string]any{
		"customerStatusCreate": map[string]any{
			"success": true,
			"customerStatus": map[string]any{
				"id":          "s1",
				"name":        "active",
				"displayName": "Active",
				"color":       "#00ff00",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	status, err := c.CreateCustomerStatus(map[string]any{"name": "active", "color": "#00ff00"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.DisplayName != "Active" {
		t.Errorf("expected Active, got %q", status.DisplayName)
	}
}

func TestUpdateCustomerStatus(t *testing.T) {
	responseData := map[string]any{
		"customerStatusUpdate": map[string]any{
			"success": true,
			"customerStatus": map[string]any{
				"id":          "s1",
				"name":        "inactive",
				"displayName": "Inactive",
				"color":       "#ff0000",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	status, err := c.UpdateCustomerStatus("s1", map[string]any{"displayName": "Inactive"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.DisplayName != "Inactive" {
		t.Errorf("expected Inactive, got %q", status.DisplayName)
	}
}

func TestDeleteCustomerStatus(t *testing.T) {
	responseData := map[string]any{
		"customerStatusDelete": map[string]any{
			"success": true,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	if err := c.DeleteCustomerStatus("s1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateCustomerTier(t *testing.T) {
	responseData := map[string]any{
		"customerTierCreate": map[string]any{
			"success": true,
			"customerTier": map[string]any{
				"id":          "t1",
				"name":        "enterprise",
				"displayName": "Enterprise",
				"color":       "#gold",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	tier, err := c.CreateCustomerTier(map[string]any{"name": "enterprise", "color": "#gold"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tier.DisplayName != "Enterprise" {
		t.Errorf("expected Enterprise, got %q", tier.DisplayName)
	}
}

func TestUpdateCustomerTier(t *testing.T) {
	responseData := map[string]any{
		"customerTierUpdate": map[string]any{
			"success": true,
			"customerTier": map[string]any{
				"id":          "t1",
				"name":        "starter",
				"displayName": "Starter",
				"color":       "#silver",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	tier, err := c.UpdateCustomerTier("t1", map[string]any{"displayName": "Starter"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tier.DisplayName != "Starter" {
		t.Errorf("expected Starter, got %q", tier.DisplayName)
	}
}

func TestDeleteCustomerTier(t *testing.T) {
	responseData := map[string]any{
		"customerTierDelete": map[string]any{
			"success": true,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	if err := c.DeleteCustomerTier("t1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
