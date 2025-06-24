#pragma once

#include<string>
#include<crow.h>
#include "searcher/Searcher.hpp"

crow::response HandleSuggests(const std::string& request, Searcher* searcher);