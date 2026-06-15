param(
    [int]$TicketTypeId  = 15,
    [string]$DbUser     = "postgres",
    [string]$DbPassword = "123",
    [string]$DbName     = "ticketing",
    [string]$ServerUrl  = "http://localhost:8080"
)

$env:PGPASSWORD = $DbPassword

# 1. Check the server is up
Write-Host ""
Write-Host "[1/4] Checking server is running at $ServerUrl ..."
$serverUp = $false
try {
    $resp = Invoke-WebRequest -Uri "$ServerUrl/api/events" -UseBasicParsing -TimeoutSec 3
    $serverUp = $true
} catch {
    $serverUp = $false
}
if (-not $serverUp) {
    Write-Host "      FAIL -- server not reachable. Run 'make dev' in another terminal first." -ForegroundColor Red
    exit 1
}
Write-Host "      OK"

# 2. Reset PostgreSQL sold_quantity
Write-Host ""
Write-Host "[2/4] Resetting sold_quantity for ticket_type id=$TicketTypeId ..."
psql -U $DbUser -h localhost -d $DbName -c "UPDATE ticket_types SET sold_quantity = 0 WHERE id = $TicketTypeId"
if ($LASTEXITCODE -ne 0) {
    Write-Host "      FAIL -- psql error. Check DB is running and credentials are correct." -ForegroundColor Red
    exit 1
}
psql -U $DbUser -h localhost -d $DbName -c "SELECT id, name, total_quantity, sold_quantity FROM ticket_types WHERE id = $TicketTypeId"

# 3. Reset Redis holds
Write-Host ""
Write-Host "[3/4] Clearing Redis holds for ticket_holds:$TicketTypeId ..."
$redisOut = redis-cli DEL "ticket_holds:$TicketTypeId"
if ($redisOut -eq "1") {
    Write-Host "      Key deleted."
} else {
    Write-Host "      Key was already empty (nothing to clear)."
}

# 4. Run k6
Write-Host ""
Write-Host "[4/4] Running k6 spike test (TICKET_TYPE_ID=$TicketTypeId) ..."
Write-Host ""
k6 run --env TICKET_TYPE_ID=$TicketTypeId test_tickets.js

# Post-run state check
Write-Host ""
Write-Host "--- Post-run DB state ---"
psql -U $DbUser -h localhost -d $DbName -c "SELECT id, name, total_quantity, sold_quantity, total_quantity - sold_quantity AS remaining FROM ticket_types WHERE id = $TicketTypeId"

Write-Host "--- Post-run Redis state ---"
$held = redis-cli HGET "ticket_holds:$TicketTypeId" held_quantity
if (-not $held) { $held = "0" }
Write-Host "ticket_holds:${TicketTypeId}  held_quantity = $held"
Write-Host "(holds release automatically after 15 min TTL)"
Write-Host ""
