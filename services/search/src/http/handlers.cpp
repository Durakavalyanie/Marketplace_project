#include "http/handlers.hpp"
#include <crow.h>


crow::response HandleSuggests(const crow::request& req, crow::response& res, std::string& request, Searcher* searcher) {
    crow::json::wvalue response;
    response["suggests"] = searcher->SearchSuggests(request);
    return crow::response{response};
}