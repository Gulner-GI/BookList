BASE_URL="http://localhost:8080/todos"

if [ "$#" -ne 1 ]; then
  echo "Использование: $0 <id>"
  exit 1
fi

id="$1"

curl -s -X DELETE "$BASE_URL/$id"