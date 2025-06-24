FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moscow

# Установка необходимых зависимостей
RUN apt update && apt install -y \
    git \
    cmake \
    g++ \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Копируем содержимое проекта и vcpkg в контейнер
WORKDIR /app
COPY . .

# Монтируем vcpkg из хоста
ENV VCPKG_ROOT=/vcpkg
ENV CMAKE_TOOLCHAIN_FILE=$VCPKG_ROOT/scripts/buildsystems/vcpkg.cmake

RUN ls -l /vcpkg

# Создаем папку для сборки
RUN mkdir -p build
WORKDIR /app/build

# Сборка проекта с использованием vcpkg
RUN rm -rf CMakeCache.txt && cmake .. -DCMAKE_TOOLCHAIN_FILE=$VCPKG_ROOT/scripts/buildsystems/vcpkg.cmake && make -j$(nproc)
# Запуск приложения
CMD ["./my_service"]