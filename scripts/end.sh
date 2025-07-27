if [ -f server.pid ]; then
  kill $(cat server.pid)
  rm server.pid
  echo "Сервер остановлен"
else
  echo "PID файл не найден"
fi