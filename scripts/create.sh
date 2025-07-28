BASE_URL="http://localhost:8080/books"

if [ "$#" -lt 4 ]; then
  echo "Использование: $0 <title> <year> <genre> <status> [link]"
  exit 1
fi

title="$1"
year="$2"
genre="$3"
status="$4"
link="$5"

json_body='{
  "title": "'"$title"'", 
  "year": '"$year"', 
  "genre": "'"$genre"'", 
  "status": "'"$status"'"'

if [ -n "$link" ]; then 
  json_body+=',"link": "'"$link"'"'
fi
json_body+='}'

curl -s -X POST -H "Content-Type: application/json" \
  -d "$json_body" \
  "$BASE_URL"