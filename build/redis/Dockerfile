FROM redis:7-alpine

COPY config/redis.conf /usr/local/etc/redis/redis.conf

CMD ["sh", "-c", "exec redis-server /usr/local/etc/redis/redis.conf --requirepass \"$REDIS_PASSWORD\""]