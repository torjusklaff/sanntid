
from threading import Thread
import thread

i = 0
i_lock = thread.allocate_lock()

def someThreadFunction_1():

# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")
    global i
    for j in xrange(1,1000000):
    	i_lock.acquire()
    	i+=1
    	i_lock.release()

def someThreadFunction_2():
    
    global i
    for j in xrange(1,1000000):
    	i_lock.acquire()
    	i-=1
    	i_lock.release()


def main():
    someThread_1 = Thread(target = someThreadFunction_1, args = (),)
    someThread_2 = Thread(target = someThreadFunction_2, args = (),)
    someThread_1.start()
    someThread_2.start()
    someThread_1.join()
    someThread_2.join()
    print(i)


main()
