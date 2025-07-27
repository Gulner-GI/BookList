go run . &
echo $! > server.pid
echo "Сервер запущен с PID $(cat server.pid)"