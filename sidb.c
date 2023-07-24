#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
    char *buffer;
    size_t buffer_length;
    size_t input_length;
} InputBuffer;

InputBuffer *new_input_buffer() {
    InputBuffer *buffer = (InputBuffer *) malloc(sizeof(InputBuffer));
    buffer->buffer = NULL;
    buffer->buffer_length = 0;
    buffer->input_length = 0;
    return buffer;
}

void print_prompt() { printf("db > "); }

void read_input(InputBuffer *input_buffer) {
    ssize_t bytes_read = getline(&(input_buffer->buffer), &(input_buffer->buffer_length), stdin);
    if (bytes_read <= 0) {
        printf("Error reading input\n");
        exit(EXIT_FAILURE);
    }
    input_buffer->input_length = bytes_read - 1; // remove newline symbol
    input_buffer->buffer[bytes_read - 1] = 0;
}

void close_buffer(InputBuffer* input_buffer) {
    free(input_buffer->buffer);
    free(input_buffer);
}

int main(int argc, char *argv[]) {
    InputBuffer *input_buffer = new_input_buffer();
    while (1) {
        print_prompt();
        read_input(input_buffer);

        if (strcmp(input_buffer->buffer, ".exit") == 0) {
            close_buffer(input_buffer);
            exit(EXIT_SUCCESS);
        } else {
            printf("Unrecognized command %s\n", input_buffer->buffer);
        }
    }
}
