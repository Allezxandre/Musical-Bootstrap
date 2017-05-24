go get # Fetch all dependencies
# Run
eval $(cat .env | sed 's/^/export /') && go build . && ./musical-bootstrap "$@"
