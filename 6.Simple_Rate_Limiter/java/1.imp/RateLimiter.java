package com.example.ratelimiter;

import redis.clients.jedis.Jedis;
import redis.clients.jedis.params.ZAddParams;

import javax.servlet.http.HttpServletResponse;
import java.time.Instant;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class RateLimiter {

    private static final String REDIS_HOST = "localhost";
    private static final int REDIS_PORT = 6379;
    private static final long DEFAULT_WINDOW_SECONDS =  24 * 60 * 60;

    private final Jedis jedis;
    private final Map<String, EndpointConfig> endpointConfigs;

    public RateLimiter() {
        this.jedis = new Jedis(REDIS_HOST, REDIS_PORT);
        this.endpointConfigs = new ConcurrentHashMap<>();
    }

    public void configure(String endpoint, int maxCredits, int creditCost) {
        configure(endpoint, maxCredits, DEFAULT_WINDOW_SECONDS, creditCost);
    }

    public void configure(String endpoint, int maxCredits, long windowSeconds, int creditCost) {
        endpointConfigs.put(endpoint, new EndpointConfig(maxCredits, windowSeconds, creditCost));
    }

    public boolean allow(String endpoint, String userId, HttpServletResponse response) {
        EndpointConfig config = endpointConfigs.get(endpoint);
        if (config == null) throw new IllegalArgumentException("Endpoint not configured");

        String redisKey = "api_credits:" + endpoint + ":" + userId;
        String resetKey = "reset_time:" + endpoint + ":" + userId;

        long currentTime = Instant.now().getEpochSecond();
        long windowStartTime = currentTime - config.windowSeconds;

        refreshDailyCredits(resetKey, redisKey, config, currentTime);

        jedis.zremrangeByScore(redisKey, 0, windowStartTime);
        long usedCredits = jedis.zcard(redisKey);

        if (usedCredits + config.creditCost > config.maxCredits) {
            setRateLimitHeaders(response, config.maxCredits, usedCredits, windowStartTime + config.windowSeconds);
            response.setStatus(429);
            return false;
        }

        for (int i = 0; i < config.creditCost; i++) {
            jedis.zadd(redisKey, currentTime, currentTime + "-" + i, ZAddParams.zAddParams().nx());
        }

        jedis.expire(redisKey, (int) config.windowSeconds);
        setRateLimitHeaders(response, config.maxCredits, usedCredits + config.creditCost, windowStartTime + config.windowSeconds);
        return true;
    }

    private void refreshDailyCredits(String resetKey, String redisKey, EndpointConfig config, long currentTime) {
        String lastResetTimeStr = jedis.get(resetKey);

        if (lastResetTimeStr == null || isNextDay(Long.parseLong(lastResetTimeStr), currentTime)) {
            jedis.del(redisKey);
            jedis.set(resetKey, String.valueOf(currentTime));
        }
    }

    private boolean isNextDay(long lastResetTime, long currentTime) {
        long oneDayInSeconds = 24 * 60 * 60;
        return currentTime - lastResetTime >= oneDayInSeconds;
    }

    private void setRateLimitHeaders(HttpServletResponse response, int maxCredits, long usedCredits, long resetTime) {
        response.setHeader("X-RateLimit-Limit", String.valueOf(maxCredits));
        response.setHeader("X-RateLimit-Remaining", String.valueOf(Math.max(0, maxCredits - usedCredits)));
        response.setHeader("X-RateLimit-Reset", String.valueOf(resetTime));
    }

    private static class EndpointConfig {
        private final int maxCredits;
        private final long windowSeconds;
        private final int creditCost;

        public EndpointConfig(int maxCredits, long windowSeconds, int creditCost) {
            this.maxCredits = maxCredits;
            this.windowSeconds = windowSeconds;
            this.creditCost = creditCost;
        }
    }
}


/**
 * 
 * 
Yes, the use of currentTime as the score in the sorted set is crucial for implementing the sliding window counter algorithm. Here's why:

Sliding Window Counter Algorithm
The sliding window counter algorithm is a rate-limiting technique that allows you to track the number of requests over a rolling time window. This is achieved by storing timestamps of requests and periodically cleaning up old entries that fall outside the window.

How It Works in This Code
Storing Timestamps:
for (int i = 0; i < config.creditCost; i++) {
    jedis.zadd(redisKey, currentTime, currentTime + "-" + i, ZAddParams.zAddParams().nx());
}

Each request's timestamp (currentTime) is added to the sorted set with a unique value (currentTime + "-" + i).
The score (timestamp) allows you to track when each credit was used.
Removing Old Entries:

jedis.zremrangeByScore(redisKey, 0, windowStartTime);

This removes all entries with a score (timestamp) less than or equal to windowStartTime, effectively removing credits that are outside the current window.
Counting Used Credits:
long usedCredits = jedis.zcard(redisKey);

This counts the number of entries in the sorted set, representing the number of credits used within the current window.
Why currentTime is Crucial
Accurate Tracking: Using the current timestamp as the score ensures that each credit is accurately tracked within the time window.
Efficient Cleanup: The sorted set allows efficient removal of old entries that fall outside the window, ensuring that only relevant credits are counted.
Sliding Window: By continuously adding new timestamps and removing old ones, the algorithm maintains a sliding window of used credits, providing a more flexible and accurate rate-limiting mechanism compared to fixed window counters.
Summary
Using currentTime as the score in the sorted set is crucial for implementing the sliding window counter algorithm. It allows accurate tracking of request timestamps, efficient cleanup of old entries, and maintains a rolling time window for rate limiting.
 * 
 */