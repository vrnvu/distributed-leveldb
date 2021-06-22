#include <stdio.h>
#include <stdbool.h>
#include <unistd.h>


int main() {
	unsigned int second = 1000000;
	while(true) {
		printf("running fake\n");
		usleep(second * 2);
	}
}