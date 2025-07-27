BASE_URL="http://localhost:8080/todos"

if [ "$#" -ne 2 ]; then
  echo "Использование: $0 <task> <done>"
  exit 1
fi

task="$1"
done="$2"

curl -s -X POST -H "Content-Type: application/json" \
  -d '{"task": "'"$task"'", "done": '"$done"'}' \
  "$BASE_URL"