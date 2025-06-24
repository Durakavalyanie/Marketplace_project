cmake_minimum_required(VERSION 3.10)

# Указываем vcpkg toolchain
set(CMAKE_TOOLCHAIN_FILE $ENV{VCPKG_ROOT}/scripts/buildsystems/vcpkg.cmake CACHE STRING "vcpkg toolchain file")

project(MyService)

# Указываем стандарт C++
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Рекурсивно собираем все файлы .cpp из папки src
file(GLOB_RECURSE SOURCES ${CMAKE_SOURCE_DIR}/src/*.cpp)

# Создаем исполняемый файл
add_executable(my_service ${SOURCES})

# Указываем директорию с заголовочными файлами
target_include_directories(my_service PRIVATE ${CMAKE_SOURCE_DIR}/include)

# Линковка с зависимостями
find_package(Crow REQUIRED)
find_package(CppKafka REQUIRED)
find_package(gRPC REQUIRED)
find_package(redis++ CONFIG REQUIRED)

# Линковка с библиотеками
target_link_libraries(my_service PRIVATE Crow::Crow)
target_link_libraries(my_service PRIVATE CppKafka::cppkafka)
target_link_libraries(my_service PRIVATE gRPC::grpc++)
target_link_libraries(my_service PRIVATE redis++::redis++)