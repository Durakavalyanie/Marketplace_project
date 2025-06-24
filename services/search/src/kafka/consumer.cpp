#include <cppkafka/consumer.h>
#include <iostream>

void kafka_listener() {
    cppkafka::Configuration config = {
        {"metadata.broker.list", "localhost:9092"},
        {"group.id", "my_group"},
        {"auto.offset.reset", "earliest"}
    };

    cppkafka::Consumer consumer(config);
    consumer.subscribe({"my_topic"});

    while (true) {
        auto msg = consumer.poll();
        if (msg && !msg.get_error()) {
            std::cout << "Kafka received: " << msg.get_payload() << std::endl;
        }
    }
}