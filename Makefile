BIN_DIR:=./bin
SRC_DIR=.

BINS:=$(BIN_DIR)/kether
MAIN_SRCS:=$(SRC_DIR)/init.go $(SRC_DIR)/main.go

.PHONY:all clean kether

all:kether

clean:
	rm -rf $(BIN_DIR)

kether:
	mkdir -p $(BIN_DIR)
	go build -o $(BINS) $(MAIN_SRCS)