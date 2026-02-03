$ErrorActionPreference = "Stop"
$baseUrl = "http://localhost:8080"

# 1. Login as Admin
$adminUser = "admin_test_v3"
$password = "password123"

Write-Host "Logging in..."
try {
    $loginBody = @{
        username = $adminUser
        password = $password
    } | ConvertTo-Json
    
    $tokenResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $tokenResponse.data.token
    $headers = @{
        Authorization = "Bearer $token"
    }
} catch {
    Write-Host "Login failed. Trying to register..."
    $registerBody = @{
        username = $adminUser
        password = $password
        role = "admin"
        api_key = "secret_key_v3"
    } | ConvertTo-Json
    Invoke-RestMethod -Uri "$baseUrl/auth/register" -Method Post -Body $registerBody -ContentType "application/json"
    
    # Login again
    $tokenResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $tokenResponse.data.token
    $headers = @{
        Authorization = "Bearer $token"
    }
}

# 2. Create a Question
Write-Host "Creating a base Question..."
$questionBody = @{
    title = "Test Question for Link"
    content = "What is the capital of France?"
    answer = "PARIS"
} | ConvertTo-Json

$qResponse = Invoke-RestMethod -Uri "$baseUrl/questions" -Method Post -Body $questionBody -ContentType "application/json" -Headers $headers
$questionId = $qResponse.data.id
Write-Host "Created Question ID: $questionId"

# 3. Create a Level (if needed, but we can use existing one or create new)
# We need a level_id for the crossword question
Write-Host "Getting a level..."
$levelsResponse = Invoke-RestMethod -Uri "$baseUrl/crosswords/levels" -Method Get -Headers $headers
$levelId = 0
if ($levelsResponse.data.Count -gt 0) {
    $levelId = $levelsResponse.data[0].id
} else {
    Write-Host "Creating a new level..."
    $levelBody = @{ name = "Test Level" } | ConvertTo-Json
    $lResponse = Invoke-RestMethod -Uri "$baseUrl/crosswords/levels" -Method Post -Body $levelBody -ContentType "application/json" -Headers $headers
    $levelId = $lResponse.data.id
}
Write-Host "Using Level ID: $levelId"

# 4. Create Crossword Question Linked to Question
Write-Host "Creating Linked Crossword Question..."
# Note: We are NOT sending clue and answer, expecting backend to fill them
$cqBody = @{
    level_id = $levelId
    number = 99
    isAcross = $true
    row = 5
    col = 5
    questions_id = $questionId
} | ConvertTo-Json

$cqResponse = Invoke-RestMethod -Uri "$baseUrl/crosswords/questions" -Method Post -Body $cqBody -ContentType "application/json" -Headers $headers
$createdCQ = $cqResponse.data

Write-Host "Created Crossword Question:"
Write-Host "ID: $($createdCQ.id)"
Write-Host "Clue: $($createdCQ.clue)"
Write-Host "Answer: $($createdCQ.answer)"
Write-Host "QuestionsID: $($createdCQ.questions_id)"

if ($createdCQ.clue -eq "What is the capital of France?" -and $createdCQ.answer -eq "PARIS") {
    Write-Host "SUCCESS: Clue and Answer were automatically populated!"
} else {
    Write-Host "FAILURE: Clue and Answer were NOT populated correctly."
    exit 1
}
