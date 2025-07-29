BASE_URL="http://localhost:8080/books"

ID=""
TITLE=""
YEAR=""
GENRE=""
STATUS=""
LINK=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -id)
      ID="$2"; shift 2;;
    -title)
      TITLE="$2"; shift 2;;
    -year)
      YEAR="$2"; shift 2;;
    -genre)
      GENRE="$2"; shift 2;;
    -status)
      STATUS="$2"; shift 2;;
    -link)
      LINK="$2"; shift 2;;
    *)
      echo "Неизвестный аргумент: $1"
      exit 1
      ;;
  esac
done

if [[ -z "$ID" ]]; then
  echo "Ошибка: нужно указать -id"
  exit 1
fi

if [[ -n "$YEAR" ]] && ! [[ "$YEAR" =~ ^[0-9]+$ ]]; then
  echo "Ошибка: год должен состоять только из цифр, а не '$YEAR'"
  exit 1
fi

if [[ -n "$LINK" ]] && ! [[ "$LINK" =~ ^https?://.+ ]]; then
  echo "Ошибка: ссылка должна начинаться с http:// или https://, а не '$LINK'"
  exit 1
fi

JSON="{"
[ -n "$TITLE" ]  && JSON+="\"title\":\"$TITLE\","
[ -n "$YEAR" ]   && JSON+="\"year\":$YEAR,"
[ -n "$GENRE" ]  && JSON+="\"genre\":\"$GENRE\","
[ -n "$STATUS" ] && JSON+="\"status\":$STATUS,"
[ -n "$LINK" ]   && JSON+="\"link\":\"$LINK\","
JSON="${JSON%,}}"

echo "Отправка JSON: $JSON"

curl -s -X PATCH "$BASE_URL/$ID" \
     -H "Content-Type: application/json" \
     -d "$JSON"