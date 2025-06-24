#include <crow.h>
#include <sw/redis++/redis++.h>
#include <iostream>
#include <memory>
#include <atomic>
#include <limits>
#include "http/server.hpp"

#define BATCH_SIZE = 10_000
#define NAMES_LIMIT = 1_000_000
std::atomic<std::shared_ptr<Searcher>> global_searcher;

void insert_list(sw::redis::Redis& redis, std::shared_ptr<Searcher> searcher, std::string key, size_t limit = std::numeric_limits<size_t>::max()) {
    size_t start_index = 0;
    std::vector<std::string> batch;

    while (start_index < limit) {
        size_t remaining = limit - start_index;
        size_t count = std::min(BATCH_SIZE, remaining);

        size_t end_index = start_index + count - 1;
        batch = redis.lrange(key, start_index, end_index);
        
        if (batch.empty()) {
            break;
        }

        for (const auto& item : batch) {
            searcher->InsertSuggestion(item);
        }

        if (limit - 1 - BATCH_SIZE < start_index) {
            break;
        }

        start_index += count;
    }
}

std::shared_ptr<Searcher> build_searcher_from_redis(sw::redis::Redis& redis) {
    std::shared_ptr<Searcher> searcher = std::make_shared<Searcher>();

    insert_list(redis, searcher, "item_names", NAMES_LIMIT)
    insert_list(redis, searcher, "recent_queries", NAMES_LIMIT)

    return searcher;
}

int start_server() {
    sw::redis::Redis redis("tcp://redis:6379");

    redis.rpush("item_names", {"rick owens jacket", "rick owens pants", "ccp boots"})
    redis.rpush("recent_queries", {"rick", "owens", "rick", "ccp", "ccp"})

    auto searcher = build_searcher_from_redis(redis);

    global_searcher.store(searcher);

    crow::SimpleApp app;

    CROW_ROUTE(app, "/suggest").methods(crow::HTTPMethod::Get)
    ([](const crow::request& req) {
        auto query = req.url_params.get("query");
        if (!query) {
            return crow::response(400, "Missing query parameter 'query'");
        }

        auto searcher = global_searcher.load();
        auto suggestions = searcher->GetSuggestions(std::string(query));

        crow::json::wvalue json;
        json["suggestions"] = crow::json::wvalue::list(suggestions.begin(), suggestions.end());
        return crow::response(json);
    });

    app.port(8081).multithreaded().run();
}