BASE_URL="http://localhost:8080/todos"

if [ -z "$1" ]; then
  echo "Использование: $0 <id> [task] [done]"
  echo "Примеры:"
  echo "  $0 1               # Ошибка, недостаточно данных"
  echo "  $0 1 'Новое дело'  # Обновит только task"
  echo "  $0 1 '' true       # Обновит только done"
  echo "  $0 1 'Новое дело' true  # Обновит оба"
  exit 1
fi

id="$1"
task="$2"
done="$3"

json="{"
first=true

if [ -n "$task" ]; then
  json="$json\"task\": \"$task\""
  first=false
fi

if [ -n "$done" ]; then
  if [ "$first" = false ]; then
    json="$json, "
  fi
  json="$json\"done\": $done"
fi

json="$json}"

curl -s -X PATCH -H "Content-Type: application/json" \
  -d "$json" \
  "$BASE_URL/$id"