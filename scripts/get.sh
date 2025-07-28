BASE_URL="http://localhost:8080/books"

if [ -z "$1" ]; then
  curl -s -X GET "$BASE_URL"
else
  id="$1"
  curl -s -X GET "$BASE_URL?id=$id"
fi