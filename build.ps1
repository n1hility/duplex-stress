go clean
$ENV:GOOS="windows"; go build
$ENV:GOOS="linux"; go build
