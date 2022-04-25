// g++ c.cpp -o c - компилируем код (-o c не обязательный параметр означает создать файл с название c.exe - по умолчанию создастся файл с именем a.exe)
// ./c            - запускаем программу

#include <iostream>
#include "factorial.h"
// создаем 2 файла factorial.h с объявлением ф-ции и factorial.cpp с самой функцией
// для компиляции 2х файлов: g++ c.cpp factorial.cpp -o c

int main()
{
    int x = 5;
    int y = x * 5;
    std::cout << "Привет!\n";
    std::cout << "x = "<< x << "\n";
    std::cout << "y = "<< y << "\n";

    int i = 0;
//    std::cin >> i;
//    printf("Iv = %d\n", i);

    while(++i < 10)
      printf("f(%d) = %d\n", i, factorial(i));

    return 0;
}
