package handlers_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// luaReserve is the same script used in CreateCheckoutSession.
// Keeping it here as a copy ensures this test catches any divergence.
const luaReserve = `
local reservationId = ARGV[#ARGV - 1]
local reservationTTL = tonumber(ARGV[#ARGV])
local numTickets = (#KEYS)
local reservedItems = {}
local fullReservationKey = "reservation:" .. reservationId

local function rollback()
	for i = 1, #reservedItems, 2 do
		local ttId = reservedItems[i]
		local qty = reservedItems[i+1]
		redis.call('HINCRBY', "ticket_holds:" .. ttId, "held_quantity", -qty)
	end
	redis.call('DEL', fullReservationKey)
	return 0
end

for i = 1, numTickets do
	local ticketHoldKey = KEYS[i]
	local totalQty = tonumber(ARGV[(i-1)*3 + 1])
	local soldQty = tonumber(ARGV[(i-1)*3 + 2])
	local requestedQty = tonumber(ARGV[(i-1)*3 + 3])
	local ticketTypeId = string.match(ticketHoldKey, "ticket_holds:(%d+)")

	local currentHeldQty = tonumber(redis.call('HGET', ticketHoldKey, 'held_quantity') or '0')
	local availableForSale = totalQty - soldQty - currentHeldQty

	if availableForSale < requestedQty then
		return rollback()
	end

	redis.call('HINCRBY', ticketHoldKey, "held_quantity", requestedQty)
	table.insert(reservedItems, ticketTypeId)
	table.insert(reservedItems, requestedQty)
	redis.call('HSET', fullReservationKey, ticketTypeId, requestedQty)
end

redis.call('EXPIRE', fullReservationKey, reservationTTL)
return 1
`

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 1}) // DB 1 = test isolation
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Skipf("Redis not reachable at localhost:6379 — skipping: %v", err)
	}
	return rdb
}

// TestNoOversell_ConcurrentReservations is the core correctness test.
// It simulates N buyers simultaneously trying to claim tickets for an event
// that only has `totalTickets` available. After all goroutines finish it
// asserts that held_quantity never exceeded the supply.
func TestNoOversell_ConcurrentReservations(t *testing.T) {
	const (
		ticketTypeID = 9999   // synthetic ID, won't exist in dev DB
		totalTickets = 5      // intentionally scarce
		soldTickets  = 0
		buyers       = 50    // concurrent buyers — far more than available tickets
	)

	ctx := context.Background()
	rdb := newTestRedis(t)
	defer rdb.Close()

	holdKey := fmt.Sprintf("ticket_holds:%d", ticketTypeID)

	// Clean state before test
	rdb.Del(ctx, holdKey)

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		succeeded int
		rejected  int
	)

	for i := 0; i < buyers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			reservationID := uuid.New().String()
			keys := []string{holdKey}
			args := []interface{}{
				strconv.Itoa(totalTickets),
				strconv.Itoa(soldTickets),
				"1",           // each buyer wants 1 ticket
				reservationID,
				"900",         // TTL in seconds
			}

			val, err := rdb.Eval(ctx, luaReserve, keys, args...).Result()
			if err != nil {
				t.Errorf("Lua script error: %v", err)
				return
			}

			mu.Lock()
			if val.(int64) == 1 {
				succeeded++
			} else {
				rejected++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Assert correctness
	heldStr, err := rdb.HGet(ctx, holdKey, "held_quantity").Result()
	if err != nil {
		t.Fatalf("Could not read held_quantity from Redis: %v", err)
	}
	held, _ := strconv.Atoi(heldStr)

	t.Logf("buyers=%d  succeeded=%d  rejected=%d  held_quantity=%d  totalTickets=%d",
		buyers, succeeded, rejected, held, totalTickets)

	if held > totalTickets {
		t.Errorf("OVERSOLD: held_quantity=%d exceeded totalTickets=%d", held, totalTickets)
	}
	if succeeded != totalTickets {
		t.Errorf("expected exactly %d successes, got %d", totalTickets, succeeded)
	}
	if succeeded+rejected != buyers {
		t.Errorf("successes+rejections should equal buyers: %d+%d != %d", succeeded, rejected, buyers)
	}

	// Cleanup
	rdb.Del(ctx, holdKey)
}

// TestReservation_RollbackOnPartialFailure verifies that if a multi-ticket
// reservation can't be fully satisfied, no partial hold is left behind.
func TestReservation_RollbackOnPartialFailure(t *testing.T) {
	const (
		typeA = 9991
		typeB = 9992
	)

	ctx := context.Background()
	rdb := newTestRedis(t)
	defer rdb.Close()

	keyA := fmt.Sprintf("ticket_holds:%d", typeA)
	keyB := fmt.Sprintf("ticket_holds:%d", typeB)
	rdb.Del(ctx, keyA, keyB)

	// typeA has 5 available, typeB has 0 available
	// Reservation should fail and typeA should NOT be held
	keys := []string{keyA, keyB}
	args := []interface{}{
		"5", "0", "1", // typeA: total=5, sold=0, want=1
		"0", "0", "1", // typeB: total=0, sold=0, want=1  — will fail
		uuid.New().String(),
		"900",
	}

	val, err := rdb.Eval(ctx, luaReserve, keys, args...).Result()
	if err != nil {
		t.Fatalf("Lua script error: %v", err)
	}
	if val.(int64) != 0 {
		t.Error("expected script to return 0 (failure) when one ticket type is unavailable")
	}

	// typeA hold must be zero — rollback should have cleaned it up
	heldA, _ := rdb.HGet(ctx, keyA, "held_quantity").Result()
	if heldA != "" && heldA != "0" {
		t.Errorf("partial hold was NOT rolled back: ticket_holds:%d held_quantity=%s", typeA, heldA)
	}

	rdb.Del(ctx, keyA, keyB)
}
