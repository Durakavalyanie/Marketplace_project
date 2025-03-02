#pragma once

#include<algorithm>
#include<string>
#include<vector>
#include<unordered_map>

template<size_t N = 8>
class RadixTrie {
  public:

  RadixTrie() : root_(new Node()) { root_->edge = ""; }

  void Insert(const std::string& word) { InsertAndReturn(word); }

  std::vector<std::string> Search(const std::string& word) {
    auto[found, node] = SearchFromNode(root_, word);
    std::vector<std::string> ans;
    for (const auto& elem : node->suggests) {
      ans.push_back(elem.second);
    }
    return ans;
  }

  private:
    struct Node;

    Node* root_;

    std::pair<bool, Node*> InsertAndReturn(const std::string& request) {
      Node* node = root_;
      size_t i = 0;
      std::vector<Node*> path = {node};
      size_t suff_size;
      while(i < request.size()) {
        suff_size = request.size() - i;
        auto next_iter = node->children.find(request[i]);
        if (next_iter == node->children.end()) {
          Node* new_node = new Node();
          new_node->edge = request.substr(i, suff_size);
          new_node->rate = 1;
          node->children[request[i]] = new_node;
          UpdatePathSuggests(path, new_node->edge, new_node->rate);
          return {false, new_node};
        }

        Node* next = (*next_iter).second;
  
        if (next->edge.size() == suff_size) {
          if (request.compare(i, next->edge.size(), next->edge) == 0) {
            ++next->rate;
            UpdatePathSuggests(path, request, next->rate);
            return {true, next};
          }
        }

        size_t com_chars = CountCommonChars(next->edge, request, 0, i);

        if (next->edge.size() <= suff_size) {
          if (com_chars < next->edge.size()) {
            node = SplitNode(node, next, request.substr(i, suff_size), com_chars);
          } else {
            node = next;
          }
        } else {
          node = SeparateNode(node, next, request.substr(i, suff_size), com_chars);
          UpdatePathSuggests(path, request, node->rate);
          return {false, node};
        }

        i += com_chars;
        path.push_back(node);
      }
      return {false, nullptr};
    }

    std::pair<bool, Node*> SearchFromNode(Node* node, const std::string& request) {
      size_t i = 0;
      size_t suff_size;
      while(i < request.size()) {
        suff_size = request.size() - i;
        auto next_iter = node->children.find(request[i]);
        if (next_iter == node->children.end()) {
          return {false, node};
        }

        Node* next = (*next_iter).second;
  
        if (next->edge.size() == suff_size) {
          if (request.compare(i, next->edge.size(), next->edge) == 0) {
            return {true, next};
          }
        }

        size_t com_chars = CountCommonChars(next->edge, request, 0, i);

        if (next->edge.size() <= suff_size) {
          if (com_chars < next->edge.size()) {
            return {false, node};
          }
          node = next;
        } else {
          return {false, node};
        }

        i += com_chars;
      }
      return {false, nullptr};
    }

    Node* SplitNode(Node* parent, Node* child, const std::string& suff, size_t com_chars) {
      Node* common = new Node();
      Node* for_suffix = new Node();
      common->edge = suff.substr(0, com_chars);
      child->edge = child->edge.substr(com_chars, child->edge.size() - com_chars);
      for_suffix->edge = suff.substr(com_chars, suff.size() - com_chars);

      parent->children[common->edge[0]] = common;
      common->children[child->edge[0]] = child;
      common->children[for_suffix->edge[0]] = for_suffix;

      for_suffix->rate = 0;

      for (auto& elem : child->suggests) {
        common->suggests.push_back({elem.first, child->edge + elem.second});
      }

      UpdateSuggests(common->suggests, child->edge, child->rate);

      return common;
    }

    Node* SeparateNode(Node* parent, Node* child, const std::string& suff, size_t com_chars) {
      Node* common = new Node();
      common->edge = suff.substr(0, com_chars);
      child->edge = suff.substr(com_chars, child->edge.size() - com_chars);

      parent->children[common->edge[0]] = common;
      common->children[child->edge[0]] = child;

      common->rate = 1;

      for (auto& elem : child->suggests) {
        common->suggests.push_back({elem.first, child->edge + elem.second});
      }

      UpdateSuggests(common->suggests, child->edge, child->rate);

      return common;
    }

    void UpdatePathSuggests(const std::vector<Node*>& path, const std::string& word, size_t rate) {
      std::string word_copy = word;
      for (Node* node : path) {
        word_copy.erase(0, node->edge.size());
        UpdateSuggests(node->suggests, word, rate);
      }
    }

    void UpdateSuggests(std::vector<std::pair<size_t, std::string>>& suggests, const std::string& word, size_t rate) {
      if (suggests.size() == N && suggests[N-1].first >= rate) {
        return;
      }

      bool is_presented = false;
      for(auto& elem : suggests) {
        if (elem.second == word) {
          elem.first = rate;
          is_presented = true;
          break;
        }
      }

      if (!is_presented) {
        if (suggests.size() < N) {
          suggests.push_back({rate, word});
        } else {
          suggests[N-1] = {rate, word};
        }
      }

      for(size_t i = suggests.size() - 1; i > 0; --i) {
        if (suggests[i].first > suggests[i-1].first) {
          std::swap(suggests[i], suggests[i-1]);
        } else {
          break;
        }
      } 
    }

    size_t CountCommonChars(const std::string& str1, const std::string& str2, size_t start1, size_t start2) {
      size_t i = 0;
      while((start1 + i) < str1.size() && (start2 + i) < str2.size() && str1[start1 + i] == str2[start2 + i]) {
        ++i;
      }
      return i;
    }
};

template<size_t N>
struct RadixTrie<N>::Node {
  std::unordered_map<char, Node*> children;
  std::string edge;
  std::vector<std::pair<size_t, std::string>> suggests;
  size_t rate = 0;

  Node() = default;
};