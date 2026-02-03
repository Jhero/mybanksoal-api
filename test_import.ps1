$ErrorActionPreference = "Stop"

$baseUrl = "http://localhost:8080"

# 1. Register Admin
$adminUser = "admin_test_v4"
$password = "password123"
$registerBody = @{
    username = $adminUser
    password = $password
    role = "admin"
    api_key = "secret_key_v4"
} | ConvertTo-Json

Write-Host "Registering admin..."
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/auth/register" -Method Post -Body $registerBody -ContentType "application/json"
    Write-Host "Registered admin: $adminUser"
} catch {
    Write-Host "Registration failed: $_"
    Write-Host "Trying login..."
}

# 2. Login
$loginBody = @{
    username = $adminUser
    password = $password
} | ConvertTo-Json

Write-Host "Logging in..."
$tokenResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $tokenResponse.data.token
$headers = @{
    Authorization = "Bearer $token"
}

# 3. Import Levels
$importBody = Get-Content -Path "import_levels.json" -Raw
Write-Host "Importing levels..."
try {
    Invoke-RestMethod -Uri "$baseUrl/crosswords/import" -Method Post -Body $importBody -ContentType "application/json" -Headers $headers
    Write-Host "Levels imported successfully."
} catch {
    Write-Host "Import failed: $_"
    exit 1
}

# 4. Get Levels
Write-Host "Getting levels..."
$levelsResponse = Invoke-RestMethod -Uri "$baseUrl/crosswords/levels" -Method Get -Headers $headers
$levels = $levelsResponse.data
Write-Host "Found $($levels.Count) levels."

if ($levels.Count -eq 0) {
    Write-Host "No levels found!"
    exit 1
}

$levelId = $levels[0].id
Write-Host "First level ID: $levelId"

# 5. Get Level Details (with questions)
Write-Host "Getting level details..."
$levelResponse = Invoke-RestMethod -Uri "$baseUrl/crosswords/levels/$levelId" -Method Get -Headers $headers
$level = $levelResponse.data
Write-Host "Level Name: $($level.name)"
Write-Host "Questions count: $($level.questions.Count)"

if ($level.questions.Count -eq 0) {
    Write-Host "No questions in level!"
} else {
    $q = $level.questions[0]
    Write-Host "First Question: $($q.clue) -> $($q.answer) (LevelID: $($q.level_id))"
}

# 6. Submit Level
$answers = @{}
foreach ($q in $level.questions) {
    $answers["$($q.number)"] = $q.answer
}
$submitBody = @{
    answers = $answers
} | ConvertTo-Json

Write-Host "Submitting level..."
try {
    $scoreResponse = Invoke-RestMethod -Uri "$baseUrl/crosswords/levels/$levelId/submit" -Method Post -Body $submitBody -ContentType "application/json" -Headers $headers
    Write-Host "Score: $($scoreResponse.data.score)"
} catch {
    Write-Host "Submit failed: $_"
}

# 7. Get Leaderboard
Write-Host "Getting leaderboard..."
$leaderboardResponse = Invoke-RestMethod -Uri "$baseUrl/crosswords/leaderboard" -Method Get -Headers $headers
$entries = $leaderboardResponse.data
foreach ($entry in $entries) {
    Write-Host "User: $($entry.username), Score: $($entry.total_score)"
}
