# root Dockerfile has opencv deps
# it's what is published and runs ./containers binary

# add make command to build image inside docker image
# ./build/Dockerfile has FROM of root docker image
# and go lang dev tools, used for local builds
