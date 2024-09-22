local auth_file = io.open("auth.txt", "r")
local auth = auth_file:read("*a")
auth_file:close()


wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Basic " .. auth

request = function()
   print(wrk.url)
   io.flush()
   return wrk.format(nil, "/api/v1/namespaces/_/actions/visitor-counter?blocking=true&result=true")
end

response = function(status, headers, body)
   -- Log or process the response status, headers, and body
   print("Status: " .. status)
   print("Body: " .. body)
   io.flush()

   -- You could also handle success/failure counting here:
   if status ~= 200 then
      print("Error: Received non-200 status code.")
   end
end