#pragma once

#include<searcher/RadixTrie.hpp>
#include <fstream>
#include <string>

class Searcher {
    public:
    Searcher() = default;

    void InsertSuggestion(const std::string& request) {
        suggester_.Insert(request);
    }

    std::vector<std::string> GetSuggestions(const std::string& request) {
        return suggester_.Search(request);
    }

    private:

    static const size_t N = 8;
    RadixTrie<N> suggester_;
};