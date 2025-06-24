#pragma once

#include <sw/redis++/redis++.h>
#include <memory>
#include "searcher/Searcher.hpp"

std::shared_ptr<Searcher> build_searcher_from_redis(sw::redis::Redis& redis) 

void start_http_server(Searcher* searcher);