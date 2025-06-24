#include "http/server.hpp"
#include "searcher/Searcher.hpp"

//#include "kafka/consumer.hpp"
//#include "grpc/client.hpp"
#include <thread>
#include <iostream>

int main() {
    Searcher searcher;
    searcher.InsertSuggest("aaa");
    searcher.InsertSuggest("aa");
    searcher.InsertSuggest("aa");

    std::cout << "IN MAIN\n";

    std::thread http_thread(start_http_server, &searcher);
   // std::thread kafka_thread(kafka_listener);
   // std::thread grpc_thread(grpc_request);

    http_thread.join();
   // kafka_thread.join();
   // grpc_thread.join();

    return 0;
}