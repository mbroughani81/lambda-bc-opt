request = function()
    local path = "/get"
    local body = "{\"K\": \"cnt\"}"
    local headers = {}

    -- Set headers for POST request
    headers["Content-Type"] = "application/json"

    -- Full HTTP request setup
    return wrk.format("POST", path, headers, body)
end