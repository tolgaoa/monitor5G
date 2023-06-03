cmake_minimum_required(VERSION 3.4)
project(sgxtestAPIclient)
enable_testing()

set(CMAKE_BUILD_TYPE Debug)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

find_library(PISTACHE_LIBRARY pistache)

include_directories(${PROJECT_SOURCE_DIR})

add_executable(${PROJECT_NAME} main.cpp)

target_include_directories(${PROJECT_NAME} PUBLIC ./)

target_link_libraries(${PROJECT_NAME} ${PISTACHE_LIBRARY} -lcurl)
