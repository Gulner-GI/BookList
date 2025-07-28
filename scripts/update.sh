BASE_URL="http://localhost:8080/books"

if [ -z "$1" ]; then
  echo "Использование: $0 <id> [title] [year] [genre] [status] [link]"
  echo "Пример:"
  echo "  $0 1   \"Новый заголовок\" 2022 \"Fantasy\" true \"https://example.com\""
  exit 1
fi

id="$1"
title="$2"
year="$3"
genre="$4"
status="$5"
link="$6"

json="{"
first=true

add_field() {
  key=$1
  value=$2
  quote=$3
  if [ "$first" = false ]; then
    json="$json, "
  fi
  if [ "$quote" = true ]; then
    json="$json\"$key\": \"$value\""
  else 
    json="$json\"$key\": $value"
  fi
  first=false
}

if [ -n "$title" ]; then
  add_field "title" "$title" true
fi

if [ -n "$year" ]; then
  add_field "year" "$year" false
fi

if [ -n "$genre" ]; then
  add_field "genre" "$genre" true
fi

if [ "$status" == "true" ]; then
  add_field "status" "true" false
elif [ "$status" == "false" ]; then
  add_field "status" "false" false
fi

if [ -n "$link" ]; then
  add_field "link" "$link" true
fi

json="$json}"

echo "SENT JSON: $json"

curl -s -X PATCH -H "Content-Type: application/json" \
  -d "$json" \
  "$BASE_URL/$id"