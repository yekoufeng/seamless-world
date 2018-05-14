cmake_minimum_required(VERSION 3.0.0)
project(unitypx VERSION 0.0.0)

INCLUDE_DIRECTORIES(${PROJECT_SOURCE_DIR}/)
INCLUDE_DIRECTORIES(/home/jcmiao/software/PhysX-3.4/PhysX_3.4/Include/)
INCLUDE_DIRECTORIES(/home/jcmiao/software/PhysX-3.4/PxShared/include/)

FILE(GLOB_RECURSE DIR_HEADERS "*.h")
AUX_SOURCE_DIRECTORY(. DIR_SRCS)

ADD_DEFINITIONS(-g -W -Wall -std=c++11 -DNDEBUG)

SET(LIBRARY_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/lib)

LINK_LIBRARIES(PhysX3Extensions PxFoundation_x64 PhysX3Common_x64 PhysX3Cooking_x64 PhysX3_x64)
LINK_DIRECTORIES(/usr/local/lib64)

ADD_LIBRARY(unitypx SHARED ${DIR_SRCS} ${DIR_HEADERS})
