// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t i_lock;
// Note the return type: void*
void* someThreadFunction1(){
	int j;
    for (j = 0; j < 1000000; ++j)
    {
    	pthread_mutex_lock(&i_lock);
    	i++;
    	pthread_mutex_unlock(&i_lock);
    }
    return NULL;
}

void* someThreadFunction2(){
	int j;
    for (j = 0; j < 1000000; ++j)
    {
    	
    	pthread_mutex_lock(&i_lock);
    	i--;
    	pthread_mutex_unlock(&i_lock);
    }
    return NULL;
}


int main(){
    pthread_t someThread1;
    pthread_t someThread2;
    pthread_create(&someThread1, NULL, someThreadFunction1, NULL);
    // Arguments to a thread would be passed here ---------^
    pthread_create(&someThread2, NULL, someThreadFunction2, NULL);
    
    pthread_join(someThread1, NULL);
    pthread_join(someThread2, NULL);
    printf("i: %i\n", i);
    return 0;
    
}