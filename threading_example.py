# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread

i = 0

def someThreadFunction_1():

# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")
    global i
    for j in xrange(1,1000000):
    	i+=1

def someThreadFunction_2():
    
    global i
    for j in xrange(1,1000000):
    	i-=1


def main():
    someThread_1 = Thread(target = someThreadFunction_1, args = (),)
    someThread_2 = Thread(target = someThreadFunction_2, args = (),)
    someThread_1.start()
    someThread_2.start()
    someThread_1.join()
    someThread_2.join()
    print(i)


main()